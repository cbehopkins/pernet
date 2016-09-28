package pernet

import (
	"log"
	"sync"
	"testing"
)

func TestBasicPing(t *testing.T) {

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
func TestBasicBulk(t *testing.T) {

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
	for i := range open_cons {
		err := conn.CloseBulkConn(i)
		check(err)
	}
}
func TestBasicUDP(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient()
	//defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	// In ths test we will start up a side UDP channel to check we can send and receive chunks of data
	for i := 0; i < 2; i++ {
		port_num, err := conn.NewUDPConn()
		check(err)
		open_cons[port_num] = struct{}{}
	}
	log.Println("Bulk connection opened, try to send something")

	var out_count sync.WaitGroup
	out_count.Add(2)
	for i := range open_cons {
		err := conn.SendRxBulkUdp(10, i)
		check(err)
		out_count.Done()
	}
	out_count.Wait()
	for i := range open_cons {
		err := conn.CloseUDPConn(i)
		check(err)
	}
}
func runTestNet(num_connections int, len_p uint, t *testing.B) {
	num_its := 1
	if t != nil {
		num_its = t.N
	}
	length := 1 << len_p

	//log.Println("Starting Client...")
	conn := NewClient()
	defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	data_2_send := make(map[int][]byte)
	data_received := make(map[int][]byte)

	// In ths test we will start up a side TCP channel to check we can send and receive chunks of data
	for i := 0; i < num_connections; i++ {
		port_num, err := conn.NewBulkConn()
		check(err)
		//log.Println("Open connection on port:", port_num)
		open_cons[port_num] = struct{}{}
		bufrs := make([]byte, length)
		conn.GenBulk(bufrs)
		data_2_send[port_num] = bufrs
		bufrr := make([]byte, length)
		data_received[port_num] = bufrr

	}
	//log.Println("Bulk connections opened, try to send something")

	var out_count sync.WaitGroup
	if t != nil {
		t.ResetTimer()
	}

	for i := 0; i < num_its; i++ {
		out_count.Add(num_connections)
		for port_num := range open_cons {
			go func(pn int) {
				//err := conn.SendRxBulk(length, pn)
				var sda []byte
				var dra []byte
				sda = data_2_send[pn]
				dra = data_received[pn]

				conn.TxRxBulk(sda, dra, pn)
				//check(err)
				out_count.Done()
				go conn.CheckData(sda, dra)

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

func BenchmarkNet_16_1k(b *testing.B) {
	runTestNet(16, 10, b)
}
func BenchmarkNet_8_1k(b *testing.B) {
	runTestNet(8, 10, b)
}
func BenchmarkNet_1_1k(b *testing.B) {
	runTestNet(1, 10, b)
}
func BenchmarkNet_2_1k(b *testing.B) {
	runTestNet(2, 10, b)
}
func BenchmarkNet_16_64k(b *testing.B) {
	runTestNet(16, 16, b)
}
func BenchmarkNet_16_256k(b *testing.B) {
	runTestNet(16, 18, b)
}
func BenchmarkNet_1_512k(b *testing.B) {
	runTestNet(1, 19, b)
}
func BenchmarkNet_2_512k(b *testing.B) {
	runTestNet(2, 19, b)
}
func BenchmarkNet_16_512k(b *testing.B) {
	runTestNet(16, 19, b)
}
func BenchmarkNet_1_1m(b *testing.B) {
	runTestNet(1, 20, b)
}
func BenchmarkNet_2_1m(b *testing.B) {
	runTestNet(2, 20, b)
}
