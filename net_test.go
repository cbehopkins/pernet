package pernet

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"testing"
)

func TcpLoopConn(conn net.Conn) {
	io.Copy(conn, conn)
	conn.Close()
}

func TcpLoopConnManual(conn net.Conn) {

	var err error
	buffer := make([]byte, 64)
	for err == nil {
		var cnt int
		cnt, err = conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				if strings.Contains(err.Error(), "connection reset by peer") {
					log.Println("Connection closed in a naughty way")
				} else {
					panic(err)
				}
			}
		}
		if cnt > 0 {
			log.Printf("Read %d bytes,%v\n", cnt, buffer)
			cntw, errw := conn.Write(buffer[:cnt])
			log.Println("Write Complete")
			if errw != nil {

				if strings.Contains(errw.Error(), "connection reset by peer") {
					log.Println("Naughty close")
				} else {
					panic(errw)
				}
			}
			if cntw != cnt {
				log.Fatalf("Unable to write %d, wrote %d, %v\n", cnt, cntw, buffer[:cnt])
			}
		}
	}
	log.Println("Copy finished")
}
func dataSrc(conn io.WriteCloser) {
	r := rand.New(rand.NewSource(1))
	count := 32
	data_2_send := make([]byte, count)

	for i, _ := range data_2_send {
		data_2_send[i] = byte(r.Int())
	}
	conn.Write(data_2_send)
	fmt.Println("Finished Writing")
	// The send can finish before we have finished reading (DUH)
	// Therefore we close in the reader, not the writer
}
func dataSnk(conn io.ReadCloser, closeit bool) {
	count := 32
	data_received := make([]byte, count)
	var bytes_read int
	var err error
	for bytes_read < count && err == nil {
		var rx_d int
		rx_d, err = conn.Read(data_received[bytes_read:])
		if rx_d > 0 {
			fmt.Println("Received from connection:", rx_d, bytes_read, count, data_received)
		}
		if err == nil {
		} else if err == io.EOF {
			log.Printf("Connection with client closed\n")
		} else if err == io.ErrClosedPipe {
			log.Printf("Connection with client Pipe closed\n")
		} else {
			log.Fatalf("Bulk Connection read error: %v\n", err)
		}

		bytes_read += rx_d
	}
	fmt.Println("Received back all the data")
	if count != bytes_read {
		log.Fatalf("Incorrect message length:%d %d\n%v\n", bytes_read, len(data_received), data_received)
	}
	if closeit {
		conn.Close()
	}
}

// Form a connecton to the supplied port number
func doConnTCP(prt_num int) (conn net.Conn, err error) {
	err = fmt.Errorf("Not Connected")
	for err != nil {
		// Keep dialing until it works
		conn, err = net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(prt_num))
		if err == nil {
		} else if strings.Contains(err.Error(), "connection refused") {
		} else {
			panic(err)
		}

	}
	return
}

func testListenTCP() int {
	err := fmt.Errorf("Not Dialed")
	var ln net.Listener
	prt_num := 8082
	for err != nil {
		// Keep dialing until it works
		ln, err = net.Listen("tcp", ":"+strconv.Itoa(prt_num))
		if err != nil {
			log.Printf("Listen error: %v\n", err)
			prt_num++
		}
	}

	log.Println("Ready to Listen")
	go func() {
		defer ln.Close()
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			if err.Error() == "use of closed network connection" {
				return
			}
			log.Fatalln("Accept error:\"%v\"\n", err)
		} else {
			go TcpLoopConn(conn)
		}
	}()

	return prt_num
}

func testListenUDP() int {
	conn, prt_num := startUDP(1001)
	go LoopConnUDP(conn)

	return prt_num
}
func TestBas(t *testing.T) {
	conna, connb := net.Pipe()
	go TcpLoopConnManual(connb)
	go dataSrc(conna)
	dataSnk(conna, true)
	//conna.Close()
	//connb.Close()
}
func TestPipeTCP(t *testing.T) {

	tconn, err := doConnTCP(testListenTCP())
	check(err)
	go dataSrc(tconn)
	dataSnk(tconn, true)
	tconn.Close()
}
func TestPipeUDP(t *testing.T) {
	tconn, err := doConnUDP(testListenUDP())
	check(err)
	go dataSrc(tconn)
	dataSnk(tconn, true)
	tconn.Close()
}
func TestMultiUDP(t *testing.T) {
	port_num := testListenUDP()
	tconn, err := doConnUDP(port_num)
	check(err)
	go dataSrc(tconn)
	dataSnk(tconn, false)
	go dataSrc(tconn)
	dataSnk(tconn, false)

	tconn.Close()
}
