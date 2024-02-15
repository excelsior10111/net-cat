package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatLogEntry(name string) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: ", name)
}

func FormatJoinNotification(name string) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"]: %s has joined our chat... ", name)
}

func FormatChatMessage(text string, name string) string {
	return fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+"][%s]: %s ", name, text)
}

func IsMessageValid(message string) bool {
	if len(message) > 300 {
		return false
	}
	for _, v := range message {
		if v <= 31 {
			return false
		}
	}
	return true
}

func ClearLine(line string) string {
	return "\r" + strings.Repeat(" ", len(line)) + "\r"
}

func LinuxLogoBanner() string {
	banner := fmt.Sprintln("Welcome to TCP-Chat!")
	banner += fmt.Sprintln("         _nnnn_")
	banner += fmt.Sprintln("	dGGGGMMb")
	banner += fmt.Sprintln("       @p~qp~~qMb")
	banner += fmt.Sprintln("       M|@||@) M|")
	banner += fmt.Sprintln("       @,----.JM|")
	banner += fmt.Sprintln("      JS^\\__/  qKL")
	banner += fmt.Sprintln("     dZP        qKRb")
	banner += fmt.Sprintln("    dZP          qKKb")
	banner += fmt.Sprintln("   fZP            SMMb")
	banner += fmt.Sprintln("   HZM            MMMM")
	banner += fmt.Sprintln("   FqM            MMMM")
	banner += fmt.Sprintln(" |    `.       | `' \\Zq")
	banner += fmt.Sprintln("_)      \\.___.,|     .'")
	banner += fmt.Sprintln("\\____   )MMMMMP|   .'")
	banner += fmt.Sprintln("     `-'       `--'")
	return banner
}
