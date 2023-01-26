package server

import (
	"net"
	"sync"

	"01.alem.school/git/aseitkha/net-cat/system"
)

type Server struct {
	Listener net.Listener
}

func CreateNewServer(addr string) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
	}, nil
}

func (server *Server) RunServer() {
	chBroadCast := make(chan system.BroadCastStatus)
	chMessage := make(chan system.BroadCastMessage)
	var muChat sync.Mutex
	var mu sync.Mutex
	chat := system.EstablishNewChat(&muChat)
	go chat.BroadCastRoutine(chBroadCast, chMessage)
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			return
		}
		newUserThread := system.CreateNewThread(conn, &mu, chat)

		go newUserThread.UserHandler(chBroadCast, chMessage)
	}
}
