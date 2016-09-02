package pernet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type Client struct {
	Conn  net.Conn
	bconn net.Conn
}

func NewClient() (conn Client) {
	// connect to this socket
	tconn, err := net.Dial("tcp", "127.0.0.1:8084")
	if err != nil {
		log.Printf("Dial error: %v\n", err)
	}
	conn.Conn = tconn
	return
}
func (iconn Client) CloseAll() {
	if iconn.bconn != nil {
		iconn.bconn.Close()
	}
	iconn.Conn.Close()
}
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

func (iconn *Client) NewBulkConn() (port_num int, err error) {
	bob := Message{Action: "BConn", Data: ""}
	snd_mess, err := MarshalMessage(bob)
	if err != nil {
		log.Fatal("Error marshalling", err)
		return
	}
	fmt.Fprintln(iconn.Conn, snd_mess)
	//////////
	// listen for reply on open connection
	fmt.Println("Waiting for response from Server")
	message, err := bufio.NewReader(iconn.Conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Read String error: %v\n", err)
	}
	log.Println("Received Message:", message)
	item, err := UnmarshalMessage(message)
	check(err)
	if item.Action != "Bonn" {
		err = fmt.Errorf("Bulk Connection is not Bonn:%s", message)
		return
	}
	// Now we have created a listener for it, open the bulk connection
	log.Println("Connect to the Bulk connection that was Bonn")
	port_num, err = strconv.Atoi(item.Data)
	check(err)
	iconn.bconn, err = net.Dial("tcp", "127.0.0.1:"+item.Data)
	check(err)
	return

}
func (iconn *Client) CloseBulkConn(port_num int) (err error) {
	// TBD embed port number in bconn struct?
	bob := Message{Action: "BConnClose", Data: strconv.Itoa(port_num)}
	snd_mess, err := MarshalMessage(bob)
	check(err)
	fmt.Fprintln(iconn.Conn, snd_mess)
	message, err := bufio.NewReader(iconn.Conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Read String error: %v\n", err)
	}
	log.Println("Received Message:", message)
	item, err := UnmarshalMessage(message)
	check(err)
	if item.Data != "ok" || item.Action != "BConnClose" {
		err = fmt.Errorf("Bulk Connection is not ok:%s", message)
		return
	}
	return
}
func (iconn Client) SendRxBulk(count int) error {
	data_2_send := make([]byte, count)
	fmt.Fprintln(iconn.bconn, data_2_send)
	fmt.Println("Sent Bulk Data")
	_, err := bufio.NewReader(iconn.bconn).ReadString('\n')
	fmt.Println("Received back the bulk data")
	if err != nil {
		if err == io.EOF {
			log.Printf("Connection with client closed\n")
			return nil
		}
		log.Fatal("Bulk Connection read error: %v\n", err)

	}
	return err
}
