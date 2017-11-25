package pernet

import (
	"testing"
	//"fmt"
)

func TestRely(t *testing.T) {
	conna, connb := unrelyConn()
	go TcpLoopConnManual(connb)
	go dataSrc(conna)
	dataSnk(conna, true)
	//conna.Close()
	//connb.Close()
}
