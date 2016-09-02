package pernet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Server() {

	log.Println("Launching Server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Printf("Listen error: %v\n", err)
	}

	for {
		log.Println("Ready to Listen")
		// accept connection on port
		conn, err := ln.Accept()
		log.Println("Heard something")
		if err != nil {
			log.Printf("Accept error: %v\n", err)
		} else {
			go HandleConnection(conn)
		}
	}
}
func HandlePing(item Message, conn net.Conn) {
	// A ping message simply returns with a pong
	item.Action = "Pong"
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
}
func HandleConnection(conn net.Conn) {
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
			HandlePing(item, conn)
		default:
			log.Fatal("Unknown message", message)
		}
	} // end for
} //End handle