package pernet

import (
	"log"
	"sync"
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
	open_cons := make(map[int]struct{})
	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	for i := 0; i < 16; i++ {
		port_num, err := conn.NewBulkConn()
		check(err)
		open_cons[port_num] = struct{}{}
	}
	log.Println("Bulk connection opened, try to send something")

	var out_count sync.WaitGroup
	out_count.Add(16)
	for i := range open_cons {
		//go func() {
		err := conn.SendRxBulk(1000, i)
		check(err)
		out_count.Done()
		//}()
	}
	out_count.Wait()
	for i := range open_cons {
		err := conn.CloseBulkConn(i)
		check(err)
	}

}
