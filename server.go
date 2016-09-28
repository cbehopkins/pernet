package pernet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type serverData struct {
	sync.RWMutex
	openConnections    map[int]net.Listener
	openUDPConnections map[int]*net.UDPConn
}

func newServerData() (sd *serverData) {
	sd = new(serverData)
	sd.openConnections = make(map[int]net.Listener)
	sd.openUDPConnections = make(map[int]*net.UDPConn)

	return sd
}
func Server() {
	sd := newServerData()
	log.Println("Launching Server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Printf("Listen error: %v\n", err)
	}
	sd.openConnections[8088] = ln
	for {
		log.Println("Ready to Listen")
		// accept connection on port
		conn, err := ln.Accept()
		log.Println("Heard something")
		if err != nil {
			if err.Error() == "use of closed network connection" {
				return
			}
			log.Fatalln("Accept error:\"%v\"\n", err)
		} else {
			go sd.HandleConnection(conn)
		}
	}
}
func (sd *serverData) HandlePing(item Message, conn net.Conn) {
	// A ping message simply returns with a pong
	item.Action = "Pong"
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
}

func (sd *serverData) HandleConnection(conn net.Conn) {
	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection with client closed\n")
				return
			}
			log.Printf("Connection read error: %v\n", err)
			return
		}
		log.Printf("Received message %s\n", message)
		item, err := UnmarshalMessage(message)
		check(err)
		switch item.Action {
		case "Ping":
			sd.HandlePing(item, conn)
		case "BConn":
			sd.HandleBConn(item, conn)
		case "BConnClose":
			sd.HandleBConnClose(item, conn)
		case "UDPConn":
			sd.HandleUDPConn(item, conn)
		case "UDPConnClose":
			sd.HandleUDPConnClose(item, conn)
		default:
			log.Fatal("Unknown message", message)
		}
	} // end for
} //End handle
