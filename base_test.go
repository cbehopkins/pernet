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
func TestBatch(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	//defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	for i := 0; i < 2; i++ {
		port_num, err := conn.NewBulkConn()
		check(err)
		open_cons[port_num] = struct{}{}
	}
	log.Println("Bulk connection opened, try to send something")

	var out_count sync.WaitGroup
	out_count.Add(2)
	for i := range open_cons {
		err := conn.SendRxBulk(1000, i)
		check(err)
		out_count.Done()
	}
	out_count.Wait()
	//for i := range open_cons {
	//		err := conn.CloseBulkConn(i)
	//		check(err)
	//	}

}
func runTestNet (num_connections int, len_p uint, t *testing.B) {
	num_its := 1
	if t != nil {
		num_its = t.N
	} 
	length := 1<< len_p

	log.Println("Starting Client...")
	conn := NewClient()
	defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	for i := 0; i < num_connections; i++ {
		port_num, err := conn.NewBulkConn()
		check(err)
		log.Println("Open connection on port:", port_num)
		open_cons[port_num] = struct{}{}
	}
	log.Println("Bulk connections opened, try to send something")

	var out_count sync.WaitGroup
	if t != nil {
		t.ResetTimer()
	}

	for i:=0;i<num_its;i++ {
	out_count.Add(num_connections)
	for port_num := range open_cons {
		go func(pn int ) {
			err := conn.SendRxBulk(length, pn)
			check(err)
			out_count.Done()
		}(port_num)
	}
	out_count.Wait()
	}
	for i := range open_cons {
		err := conn.CloseBulkConn(i)
		check(err)
	}
}

func TestParrallel(t *testing.T) {
	runTestNet(16, 10, nil)
}


func BenchmarkNet(b *testing.B)  {
	runTestNet(16, 10, b)
}
