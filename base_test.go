package pernet

import (
	"log"
	"testing"
)

func TestConn(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	defer conn.Close()
	//////////
	// Now send to an open connection
	err := SendPing(conn)
	check(err)
	// Send 2 pings to make sure the connection can do this
	err = SendPing(conn)
	check(err)
}

func TestParrallel(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	defer conn.Close()
	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	bconn, err := NewBulkConn(conn)
	check(err)
	defer bconn.Close()
	err = SendRxBulk(1000, conn)
	check(err)

}
