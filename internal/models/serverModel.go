package models

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"net-cat/internal/utils"
)

type ChatMessage struct {
	Name string
	Text string
}

type ChatClient struct {
	Name string
	Conn net.Conn
}

type ChatServer struct {
	Clients    []ChatClient
	ListenAddr string
	Listener   net.Listener
	Quitch     chan struct{}
	Mesgch     chan ChatMessage
	Size       int
	History    string
	Mutex      *sync.Mutex
}

func CreateChatServer(listenAddr string) *ChatServer {
	return &ChatServer{
		ListenAddr: listenAddr,
		Quitch:     make(chan struct{}),
		Mesgch:     make(chan ChatMessage),
		Clients:    []ChatClient{},
		Mutex:      &sync.Mutex{},
		Size:       0,
	}
}

func (cs *ChatServer) Launch() error {
	listener, err := net.Listen("tcp", cs.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	cs.Listener = listener

	go cs.listenForClients()

	<-cs.Quitch
	close(cs.Mesgch)
	return nil
}

func (cs *ChatServer) listenForClients() {
	for {
		conn, err := cs.Listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		cs.Size++
		if cs.Size > 10 {
			fmt.Fprintln(conn, "chat is full, try later...")
			cs.Size = 10
			continue
		} else {
			fmt.Fprintln(conn, utils.LinuxLogoBanner())
			go cs.manageClientConnection(conn)
		}
	}
}

func (cs *ChatServer) manageClientConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	fmt.Fprint(conn, "Enter your name: ")
	var chatClient ChatClient
	for scanner.Scan() {
		chatClient.Name = scanner.Text()
		chatClient.Name = strings.Trim(chatClient.Name, " ")
		err := cs.validateClientName(chatClient.Name)
		if err != nil {
			fmt.Fprint(conn, err)
			continue
		}
		break
	}

	chatClient.Conn = conn
	fmt.Fprint(chatClient.Conn, cs.retrieveChatHistory())

	cs.Clients = append(cs.Clients, chatClient)
	fmt.Fprint(conn, utils.FormatLogEntry(chatClient.Name))

	cs.Mesgch <- ChatMessage{chatClient.Name, utils.FormatJoinNotification(chatClient.Name)}

	for scanner.Scan() {
		fmt.Fprint(conn, utils.FormatLogEntry(chatClient.Name))
		text := strings.Trim(scanner.Text(), " ")
		if len(text) == 0 {
			continue
		}
		text = utils.FormatChatMessage(text, chatClient.Name)
		if !utils.IsMessageValid(text) {
			continue
		}
		cs.Mesgch <- ChatMessage{chatClient.Name, text}
	}
	cs.Size--
	cs.removeClientFromChat(chatClient)
	cs.Mesgch <- ChatMessage{chatClient.Name, fmt.Sprintf("[%s]: %s has left our Chat. ", time.Now().Format("2006-01-02 15:04:05"), chatClient.Name)}
}

func (cs *ChatServer) removeClientFromChat(chatClient ChatClient) {
	for i := 0; i < len(cs.Clients); i++ {
		if cs.Clients[i].Name == chatClient.Name {
			cs.Clients = append(cs.Clients[:i], cs.Clients[i+1:]...)
			break
		}
	}
}

func (cs *ChatServer) retrieveChatHistory() string {
	cs.Mutex.Lock()
	defer cs.Mutex.Unlock()
	return cs.History
}

func (cs *ChatServer) validateClientName(name string) error {
	errorPromptSuffix := " try again: "
	if len(name) == 0 {
		return errors.New("name can't be empty." + errorPromptSuffix)
	}
	if len(name) > 20 {
		return errors.New("length of name must contain max 20 letters." + errorPromptSuffix)
	}
	for _, v := range name {
		if !((v >= 'A' && v <= 'Z') || (v >= 'a' && v <= 'z') || v >= '0' && v <= '9') {
			return errors.New("name has not-valid characters. allowed[a-zA-Z]." + errorPromptSuffix)
		}
	}
	for _, v := range cs.Clients {
		if v.Name == name {
			return errors.New("name you entered already exists." + errorPromptSuffix)
		}
	}
	return nil
}

func (cs *ChatServer) DistributeChatMessages() {
	var chatMessage ChatMessage
	for {
		chatMessage = <-cs.Mesgch
		cs.Mutex.Lock()
		cs.History += chatMessage.Text
		cs.History += "\n"
		for i := 0; i < len(cs.Clients); i++ {
			if cs.Clients[i].Name != chatMessage.Name {
				fmt.Fprint(cs.Clients[i].Conn, utils.ClearLine(fmt.Sprintf("[%s][%s]:", time.Now().Format("2006-01-02 15:04:05"), cs.Clients[i].Name))+chatMessage.Text+fmt.Sprintf("\n[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), cs.Clients[i].Name))
			}
		}
		cs.Mutex.Unlock()
	}
}
