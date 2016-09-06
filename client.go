package pernet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
)

type Client struct {
	Conn  net.Conn
	bconn map[int]net.Conn // map from port number to connection number
	r     *rand.Rand
}

func NewClient() (conn Client) {
	// connect to this socket
	tconn, err := net.Dial("tcp", "127.0.0.1:8084")
	if err != nil {
		log.Printf("Dial error: %v\n", err)
	}
	conn.Conn = tconn
	conn.r = rand.New(rand.NewSource(1))

	return
}
func (iconn Client) CloseAll() {
	for _, v := range iconn.bconn {

		if v != nil {
			v.Close()
		}
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
	if iconn.bconn == nil {
		iconn.bconn = make(map[int]net.Conn)
	}
	iconn.bconn[port_num], err = net.Dial("tcp", "127.0.0.1:"+item.Data)
	check(err)
	return

}
func (iconn *Client) CloseBulkConn(port_num int) (err error) {
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
	iconn.bconn[port_num].Close()
	delete(iconn.bconn, port_num)
	return
}
func (iconn Client) SendRxBulk(count, port_num int) error {
	data_2_send := make([]byte, count)
	for i, _ := range data_2_send {
		data_2_send[i] = byte(iconn.r.Int())
	}
	//fmt.Fprintln(iconn.bconn[port_num], data_2_send)
	iconn.bconn[port_num].Write(data_2_send)
	iconn.bconn[port_num].Write([]byte("\n"))

	fmt.Println("Sent Bulk Data:", port_num)
	//data_received, err := bufio.NewReader(iconn.bconn[port_num]).ReadString('\n')
	data_received := make([]byte, 0, count)
	var bytes_read int
	for bytes_read < count {
		rx_d, err := iconn.bconn[port_num].Read(data_received)
		fmt.Println("Received from cponnection:", rx_d, bytes_read, count, data_received)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection with client closed\n")
				return nil
			}
			log.Fatal("Bulk Connection read error: %v\n", err)
		}
		bytes_read += rx_d
	}
	fmt.Println("Received back the bulk data:", port_num)

	if len(data_2_send) != len(data_received) {
		log.Fatalf("Incorrect message length:%d,%d\ntype=%T,%T\n%v\n%v\n", len(data_2_send), len(data_received), data_2_send, data_received, data_2_send, data_received)
	}
	for i, val := range data_2_send {
		if data_received[i] != val {
			log.Fatalf("Incorrect messages:%d\ntype=%t,%t\n%v\n%v\n", i, data_2_send, data_received, data_2_send, data_received)
		}
	}
	return nil
}
