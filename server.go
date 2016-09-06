package pernet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type serverData struct {
	sync.RWMutex
	openConnections map[int]net.Listener
}

func newServerData() (sd *serverData) {
	sd = new(serverData)
	sd.openConnections = make(map[int]net.Listener)
	return sd
}
func Server() {
	sd := newServerData()
	log.Println("Launching Server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Printf("Listen error: %v\n", err)
	}
	sd.openConnections[8084] = ln
	for {
		log.Println("Ready to Listen")
		// accept connection on port
		conn, err := ln.Accept()
		log.Println("Heard something")
		if err != nil {
			if err.Error() == "use of closed network connection" {
				return
			}
			log.Fatalln("Accept error:\"%v\"\n", err)
		} else {
			go sd.HandleConnection(conn)
		}
	}
}
func (sd *serverData) HandlePing(item Message, conn net.Conn) {
	// A ping message simply returns with a pong
	item.Action = "Pong"
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
}
func (sd *serverData) findFreePort() (i int) {
	sd.Lock()
	defer sd.Unlock()
	ok := true
	for i = 8084; i < (1<<16) && ok; {
		_, ok = sd.openConnections[i]
		if ok {
			//fmt.Println("Port already open:", i)
			i++
		}
	}
	sd.openConnections[i] = nil
	fmt.Println("Chosen Port:", i)
	return i
}
func (sd *serverData) HandleBConnClose(item Message, conn net.Conn) {
	port_to_close, err := strconv.Atoi(item.Data)
	check(err)
	//log.Println("Trying to get lock")
	sd.RLock()
	//log.Println("got Lock")
	ln, ok := sd.openConnections[port_to_close]
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
	log.Println("Closed connection port:", port_to_close)
}

func (sd *serverData) HandleBConn(item Message, conn net.Conn) {
	// Open up a new channel on specified Port
	fmt.Println("Starting Bulk connection with port:", item.Data)
	free_port := sd.findFreePort()
	prt_string := fmt.Sprintf(":%s", strconv.Itoa(free_port))
	ln, err := net.Listen("tcp", prt_string)
	sd.openConnections[free_port] = ln
	if err != nil {
		log.Printf("Listen error: %v\n", err)
	}

	go func() {
		// TBD there is no mechanism to stop this routine
		// Fix this
		for {
			log.Println("Ready to Listen on Bulk Channel")
			// accept connection on port
			conn, err := ln.Accept()
			//log.Println("Heard something on Bulk Channel")
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return
				}
				log.Fatalln("Accept error: %v\n", err)
			} else {
				go HandleBulkConnection(conn, free_port)
			}
		}
	}()

	item.Action = "Bonn"
	item.Data = strconv.Itoa(free_port) // TBD remove port from free list on close
	snd_mess, err := MarshalMessage(item)
	check(err)
	fmt.Fprintln(conn, snd_mess)
}
func HandleBulkConnection(conn net.Conn, port_num int) {
	log.Println("Port Forward started on port ", port_num)
	//io.Copy(conn, conn)
	var err error
	buffer := make([]byte, 16)
	for err == nil {
		var cnt int
		cnt, err = conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				//panic (err)
				//add in check for "connection reset by peer"
				//if strings.Contains(err.Error(), "use of closed network connection") {
			}
		}
		fmt.Printf("Read %d bytes,on port %d,%v\n", cnt, port_num, buffer)
		cntw, errw := conn.Write(buffer[:cnt])
		//check (errw)
		if errw != nil {
			if strings.Contains(errw.Error(), "connection reset by peer") {
			} else {
				panic(err)
			}
		}
		if cntw != cnt {
			log.Fatalf("Unable to write %d, wrote %d, %v\n", cnt, cntw, buffer[:cnt])
		}
	}
	log.Println("Copy finished on port ", port_num)
	// No need to close a closed connection
	//conn.Close()
}
func (sd *serverData) HandleConnection(conn net.Conn) {
	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection with client closed\n")
				return
			}
			log.Printf("Connection read error: %v\n", err)
			return
		}
		log.Printf("Received message %s\n", message)
		item, err := UnmarshalMessage(message)
		check(err)
		switch item.Action {
		case "Ping":
			sd.HandlePing(item, conn)
		case "BConn":
			sd.HandleBConn(item, conn)
		case "BConnClose":
			sd.HandleBConnClose(item, conn)
		default:
			log.Fatal("Unknown message", message)
		}
	} // end for
} //End handle
