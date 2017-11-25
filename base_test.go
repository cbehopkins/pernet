package pernet

import (
	"log"
	//	"os"
	"sync"
	"testing"
)

//var remote_target string
var remote_target = "192.168.0.28"

//func TestMain(m *testing.T) {
//	remote_target = "192.168.0.28"
//	os.Exit(m.Run())
//}

func TestBasicPing(t *testing.T) {

	log.Println("Starting Client...")

	conn := NewClient(remote_target)
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

	//log.Println("Starting Client...")

	conn := NewClient(remote_target)
	defer conn.CloseAll()
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

	conn := NewClient(remote_target)
	defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	// In ths test we will start up a side UDP channel to check we can send and receive chunks of data
	num_connections := 2
	for i := 0; i < num_connections; i++ {
		port_num, err := conn.NewUDPConn()
		check(err)
		open_cons[port_num] = struct{}{}
		log.Println("Connection on port :", port_num)
	}
	log.Println("UDP connection opened, try to send something")

	var out_count sync.WaitGroup
	out_count.Add(num_connections)
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
func runTestTcp(num_connections int, len_p uint, t *testing.B) {
	num_its := 1
	if t != nil {
		num_its = t.N
	}
	length := (1 << len_p) / num_connections

	//log.Println("Starting Client...")
	conn := NewClient(remote_target)
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
				go func(a, b []byte) {
					conn.CheckData(a, b)
					out_count.Done()
				}(sda, dra)

			}(port_num)
		}
		out_count.Wait()
	}
	for i := range open_cons {
		err := conn.CloseBulkConn(i)
		check(err)
	}
}
func runTestUdp(num_connections int, len_p uint, t *testing.B, iters int) {

	//log.Println("Starting Client...")
	length := (len_p << 1) / uint(num_connections)

	conn := NewClient(remote_target)
	defer conn.CloseAll()
	open_cons := make(map[int]struct{})
	data_2_send := make(map[int][]byte)
	data_received := make(map[int][]byte)
	// In ths test we will start up a side UDP channel to check we can send and receive chunks of data
	for i := 0; i < num_connections; i++ {
		port_num, err := conn.NewUDPConn()
		check(err)
		open_cons[port_num] = struct{}{}

		bufrs := make([]byte, length)
		conn.GenBulk(bufrs)
		data_2_send[port_num] = bufrs
		bufrr := make([]byte, length)
		data_received[port_num] = bufrr
	}
	//log.Println("Bulk connection opened, try to send something")

	var out_count sync.WaitGroup
	for i := 0; i < iters; i++ {

		out_count.Add(num_connections)
		for port_num := range open_cons {
			go func(pn int) {
				//err := conn.SendRxBulkUdp(int(length), pn)
				var sda []byte
				var dra []byte
				sda = data_2_send[pn]
				dra = data_received[pn]

				conn.GenBulk(sda)
				conn.TxRxBulkUdp(sda, dra, pn)
				go func(a, c []byte) {
					conn.CheckData(a, c)
					out_count.Done()
				}(sda, dra)
			}(port_num)
		}
		out_count.Wait()
	}
	for i := range open_cons {
		err := conn.CloseUDPConn(i)
		check(err)
	}
}
func TestParrallelTcp(t *testing.T) {
	runTestTcp(16, 10, nil)
}
func TestParrallelUdp(t *testing.T) {
	runTestUdp(16, 10, nil, 1)
}
