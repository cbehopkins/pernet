package pernet

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func SendPing(conn net.Conn) (err error) {
	//fmt.Fprintf(conn, string("Ping")+"\n")
	bob := Message{Action: "Ping"}
	snd_mess, err := MarshalMessage(bob)
	if err != nil {
		log.Fatal("Error marshalling", err)
		return
	}
	fmt.Fprintln(conn, snd_mess)
	//////////
	// listen for reply on open connection
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Read String error: %v\n", err)
	}
	//log.Println("Received Message:", message)
	item, err := UnmarshalMessage(message)
	check(err)
	if item.Action != "Pong" {
		return fmt.Errorf("Ping was not Ponged:%s", message)
	}
	return
}
func NewClient() (conn net.Conn) {
	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:8084")
	if err != nil {
		log.Printf("Dial error: %v\n", err)
	}
	return
}
