package system

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type Chat struct {
	HistoryBuffer []string
	mu            *sync.Mutex
}

type BroadCastMessage struct {
	Name string
	Msg  string
}

type BroadCastStatus struct {
	IsConnected bool
	Name        string
}

func EstablishNewChat(mu *sync.Mutex) *Chat {
	return &Chat{
		mu: mu,
	}
}

func (chat *Chat) BroadCastRoutine(chBroadCast chan BroadCastStatus, chMessage chan BroadCastMessage) {
	for {
		select {
		case status := <-chBroadCast:
			if !status.IsConnected {
				message := fmt.Sprintf("%s has left the channel", status.Name)
				chat.send(message, status.Name)
			} else {
				message := fmt.Sprintf("%s has joined the channel", status.Name)
				chat.send(message, status.Name)
			}
		case incomingMsg := <-chMessage:
			log.Println("Got a message")
			chat.send(incomingMsg.Msg, incomingMsg.Name)
		}
	}

}

func (chat *Chat) send(msg string, name string) {
	//Check for validity
	//Check for empty message
	chat.mu.Lock()
	chat.HistoryBuffer = append(chat.HistoryBuffer, msg)
	chat.mu.Unlock()

	// msgFormat := FormatMsg(name, msg) // Preparing the formatted data

	for user, Conn := range totalUsers { // Range over all users and print msg to others
		if user != name {
			fmt.Fprint(Conn, "\n"+msg)
			fmt.Fprint(Conn, "\n"+FormatMsg(user, ""))

		}
	}
}

func FormatMsg(name string, msg string) string { //Formats the message properly

	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	// fmt.Println(formattedTime)
	formattedInput := fmt.Sprintf("\n[%s][%s]:%s", formattedTime, name, msg)
	return strings.TrimSpace(formattedInput)

}
