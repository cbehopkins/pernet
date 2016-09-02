package pernet

import (
	"log"
	"testing"
)

func TestConn(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	//////////
	// Now send to an open connection
	err := SendPing(conn)
	check(err)
	err = SendPing(conn)
	check(err)
	conn.Close()
}
