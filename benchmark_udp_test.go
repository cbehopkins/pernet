package pernet

import (
	"testing"
)

func BenchmarkUdp_16_1k(b *testing.B) {
	runTestUdp(16, 10, b, b.N)
}
func BenchmarkUdp_8_1k(b *testing.B) {
	runTestUdp(8, 10, b, b.N)
}
func BenchmarkUdp_1_1k(b *testing.B) {
	runTestUdp(1, 10, b, b.N)
}
func BenchmarkUdp_2_1k(b *testing.B) {
	runTestUdp(2, 10, b, b.N)
}
func BenchmarkUdp_1_64k(b *testing.B) {
	runTestUdp(1, 16, b, b.N)
}
func BenchmarkUdp_2_64k(b *testing.B) {
	runTestUdp(2, 16, b, b.N)
}
func BenchmarkUdp_8_64k(b *testing.B) {
	runTestUdp(8, 16, b, b.N)
}
func BenchmarkUdp_16_64k(b *testing.B) {
	runTestUdp(16, 16, b, b.N)
}
func BenchmarkUdp_1_256k(b *testing.B) {
	runTestUdp(1, 18, b, b.N)
}
func BenchmarkUdp_2_256k(b *testing.B) {
	runTestUdp(2, 18, b, b.N)
}
func BenchmarkUdp_4_256k(b *testing.B) {
	runTestUdp(4, 18, b, b.N)
}
func BenchmarkUdp_8_256k(b *testing.B) {
	runTestUdp(8, 18, b, b.N)
}

func BenchmarkUdp_16_256k(b *testing.B) {
	runTestUdp(16, 18, b, b.N)
}
func BenchmarkUdp_1_512k(b *testing.B) {
	runTestUdp(1, 19, b, b.N)
}
func BenchmarkUdp_2_512k(b *testing.B) {
	runTestUdp(2, 19, b, b.N)
}
func BenchmarkUdp_16_512k(b *testing.B) {
	runTestUdp(16, 19, b, b.N)
}
func BenchmarkUdp_1_1m(b *testing.B) {
	runTestUdp(1, 20, b, b.N)
}
func BenchmarkUdp_2_1m(b *testing.B) {
	runTestUdp(2, 20, b, b.N)
}
