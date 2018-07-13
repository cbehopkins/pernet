package pernet

import (
	"strconv"
	"testing"
)

func BenchmarkUdp(b *testing.B) {
	type dg struct {
		v uint
		s string
	}
	//mWorkers := []int{1, 2, 4, 8, 16, 32, 64}
	numWorkers := []int{1}
	//zeBlocks := []uint{10, 16, 18, 20}
	sizeBlocks := []dg{
		{10, "1k"},
		{11, "2k"},
		{12, "4k"},
		{13, "8k"},
		{14, "16k"},
		{15, "32k"},
		{16, "64k"},
		{17, "128k"},
		{18, "256k"},
		{19, "512k"},
		{20, "1M"},
		{21, "2M"},
		{22, "4M"},
		{23, "8M"},
		{24, "16M"},
	}
	for _, workerCnt := range numWorkers {
		wcs := strconv.Itoa(workerCnt)
		for _, sb := range sizeBlocks {
			runStr := "W:" + wcs + "_" + "S:" + sb.s
			b.Run(runStr, func(br *testing.B) {
				runTestUdp(workerCnt, sb.v, br, br.N)
			})
		}
	}
}

//func BenchmarkUdp_1_1k(b *testing.B)
//func BenchmarkUdp_2_1k(b *testing.B) {
//	runTestUdp(1, 10, b, b.N)
//}
