package pernet

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func (sd *serverData) findFreePort() (i int) {
	sd.Lock()
	defer sd.Unlock()
	ok := true
	for i = 8088; i < (1<<16) && ok; {
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

	err := fmt.Errorf("Not Dialed")
	var ln net.Listener
	for err != nil {
		// Keep dialing until it works
		ln, err = net.Listen("tcp", ":"+strconv.Itoa(free_port))
		sd.openConnections[free_port] = ln
		if err != nil {
			log.Printf("Listen error: %v\n", err)

			free_port = sd.findFreePort()
		}
	}

	sd.openConnections[free_port] = ln

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
	//log.Println("Port Forward started on port ", port_num)
	io.Copy(conn, conn)
	//log.Println("Copy finished on port ", port_num)
	// No need to close a closed connection
	//conn.Close()
}
func HandleBulkConnectionDebug(conn net.Conn, port_num int) {
	log.Println("Port Forward started on port ", port_num)
	var err error
	buffer := make([]byte, 16)
	for err == nil {
		var cnt int
		cnt, err = conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				if strings.Contains(err.Error(), "connection reset by peer") {
				} else {
					panic(err)
				}
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
	//log.Println("Copy finished on port ", port_num)
	// No need to close a closed connection
	//conn.Close()
}
