package pernet

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func startUDP(prt_num int) (conn *net.UDPConn, port_num int) {
	err := fmt.Errorf("Not Dialed")
	for err != nil {
		addr := net.UDPAddr{
			Port: prt_num,
			IP:   net.ParseIP("127.0.0.1"),
		}
		log.Println("Trying listening on port:", prt_num)
		conn, err = net.ListenUDP("udp", &addr)
		// Keep dialing until it works
		if err != nil {
			log.Printf("UDP Listen error: %v\n", err)
			prt_num++
		}
	}
	log.Println("UDP on port established")
	return conn, prt_num
}
func LoopConnUDP(conn *net.UDPConn) {

	var err error
	var addr *net.UDPAddr
	buffer := make([]byte, 64)
	for err == nil {
		var cnt int
		cnt, addr, err = conn.ReadFromUDP(buffer)
		if err != nil {
			if err != io.EOF {
				if strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("Port connection closed")
				} else if strings.Contains(err.Error(), "connection reset by peer") {
					log.Println("Connection closed in a naughty way")
				} else {
					panic(err)
				}
			}
		}
		if cnt > 0 {
			log.Printf("Read %d bytes,%v\n", cnt, buffer)
			cntw, errw := conn.WriteToUDP(buffer[:cnt], addr)
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

// Make a udp connection to the proffered port number
func doConnUDP(prt_num int) (conn net.Conn, err error) {
	err = fmt.Errorf("Not Connected")
	for err != nil {
		// Keep dialing until it works
		log.Println("Trying to dial UDP on port", prt_num)
		conn, err = net.Dial("udp", "127.0.0.1:"+strconv.Itoa(prt_num))
		if err == nil {
		} else if strings.Contains(err.Error(), "connection refused") {
		} else {
			panic(err)
		}

	}
	fmt.Println("UDP Dial succeeded")
	return
}
func (sd *serverData) HandleUDPConn(item Message, conn net.Conn) {
	// Step 1 open a UDP port
	var udp_conn *net.UDPConn
	var prt_num int
	udp_conn, prt_num = startUDP(1001)

	// Step 2, start up the receiver on that port
	sd.Lock()
	sd.openUDPConnections[prt_num] = udp_conn
	sd.Unlock()
	go LoopConnUDP(udp_conn)

	// Step 3, Say which port we have done this on
	// A ping message simply returns with a pong
	item.Action = "UDP_Ok"
	item.Data = strconv.Itoa(prt_num) // TBD remove port from free list on close
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
}

func (sd *serverData) HandleUDPConnClose(item Message, conn net.Conn) {
	port_to_close, err := strconv.Atoi(item.Data)
	check(err)
	//log.Println("Trying to get lock")
	sd.RLock()
	//log.Println("got Lock")
	ln, ok := sd.openUDPConnections[port_to_close]
	sd.RUnlock()
	//log.Println("Status is:", ok)
	if ok && ln != nil {

		err = ln.Close()
		fmt.Println("Port closed:", port_to_close)
		check(err)
		item.Data = "ok"
	} else {
		fmt.Println("inactive connection")
		item.Data = "Inactive Connection"
	}
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
	log.Println("Closed UDP connection port:", port_to_close)
}
