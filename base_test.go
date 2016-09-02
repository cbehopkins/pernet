package pernet

import (
	"log"
	"testing"
)

func TestConn(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	defer conn.Conn.Close()
	//////////
	// Now send to an open connection
	err := SendPing(conn.Conn)
	check(err)
	// Send 2 pings to make sure the connection can do this
	err = SendPing(conn.Conn)
	check(err)
}

func TestParrallel(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	defer conn.CloseAll()
	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	err := conn.NewBulkConn()
	check(err)
	log.Println("Bulk connection opened, try to send something")

	err = conn.SendRxBulk(1000)
	check(err)

}
