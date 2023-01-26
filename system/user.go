package system

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
)

type UserThread struct {
	conn net.Conn
	name string
	chat *Chat
	mu   *sync.Mutex
}

var totalUsers = make(map[string]net.Conn, 10)

func (user *UserThread) UserHandler(chBroad chan BroadCastStatus, chMsg chan BroadCastMessage) {
	if len(totalUsers) > 10 { // Checking if the lobby is full
		user.LobbyIsFull()
		return
	}
	intro, err := Intro() // Preparing the Intro message to the user
	if err != nil {
		return
	}
	_, err = user.conn.Write([]byte(intro)) // Writing the Intro message to the user
	if err != nil {
		log.Println("Couldn't print to the connection")
		return
	}
	user.AddNewName() // adding a new user name
	user.mu.Lock()
	chBroad <- BroadCastStatus{ // BroadCaseting the status of a new User
		Name:        user.name,
		IsConnected: true,
	}
	user.mu.Unlock()
	reader := bufio.NewReader(user.conn) // Create user input reader
	for {                                // Listen for user input while connected
		fmt.Fprint(user.conn, FormatMsg(user.name, "\n"))
		text, err := reader.ReadString('\n') // Store the msg that was sent by a user
		fmt.Println(text, err)
		if err == io.EOF {
			user.mu.Lock()
			delete(totalUsers, user.name)
			chBroad <- BroadCastStatus{IsConnected: false, Name: user.name}
			user.mu.Unlock()
			log.Println("Got out")
			return
		}
		if err != nil {
			log.Println("Error while Reading User input")
			return
		}

		if text != "" {
			log.Println("Got a message")
			user.mu.Lock()
			chMsg <- BroadCastMessage{
				Name: user.name,
				Msg:  text,
			}
			user.mu.Unlock()
		}
	}
}

func Intro() (string, error) {
	fileText, err := ioutil.ReadFile("assets/logo.txt")
	if err != nil {
		log.Println("Couldn't open the welcome logo file")
		return "", err
	}
	introMsg := string(fileText)
	introMsg += "\n[ENTER YOUR NAME]: "
	return introMsg, nil
}

func CreateNewThread(conn net.Conn, mu *sync.Mutex, chat *Chat) *UserThread {
	return &UserThread{
		conn: conn,
		mu:   mu,
		chat: chat,
	}
}

func (user *UserThread) AddNewName() {
	reader := bufio.NewReader(user.conn)
	temp, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Couldn't read name")
		return
	}
	temp = strings.TrimSpace(temp)

	if _, ok := totalUsers[temp]; ok {
		fmt.Fprint(user.conn, "Specified name already exists. Please try a new name\n")
		user.AddNewName()
	}
	user.name = temp
	totalUsers[temp] = user.conn
}

func (user *UserThread) LobbyIsFull() {
	fmt.Fprintf(user.conn, "We are sorry. However, the lobby is currently at its full capacity...\nPlease try again later\n")
}
