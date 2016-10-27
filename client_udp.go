package pernet

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

func (iconn *Client) NewUDPConn() (port_num int, err error) {
	bob := Message{Action: "UDPConn", Data: ""}
	snd_mess, err := MarshalMessage(bob)
	if err != nil {
		log.Fatal("Error marshalling", err)
		return
	}

	fmt.Fprintln(iconn.Conn, snd_mess)
	//////////
	// listen for reply on open connection
	//fmt.Println("Waiting for response from Server")
	message, err := bufio.NewReader(iconn.Conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Read String error: %v\n", err)
	}
	//log.Println("Received Message:", message)
	item, err := UnmarshalMessage(message)
	check(err)
	if item.Action != "UDP_Ok" {
		err = fmt.Errorf("UDP Connection is not ok:%s", message)
		return
	}
	// Now we have created a listener for it, open the bulk connection
	port_num, err = strconv.Atoi(item.Data)
	check(err)
	if iconn.uconn == nil {
		iconn.uconn = make(map[int]net.Conn)
	}
	ra_full := iconn.Conn.RemoteAddr().String()
	ra, _, err := net.SplitHostPort(ra_full)
	check(err)
	iconn.uconn[port_num], err = doConnUDP(ra, port_num)
	check(err)
	return
}
