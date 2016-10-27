package pernet

import (
	"testing"
)

func BenchmarkTcp_16_1k(b *testing.B) {
	runTestTcp(16, 10, b)
}
func BenchmarkTcp_8_1k(b *testing.B) {
	runTestTcp(8, 10, b)
}
func BenchmarkTcp_1_1k(b *testing.B) {
	runTestTcp(1, 10, b)
}
func BenchmarkTcp_2_1k(b *testing.B) {
	runTestTcp(2, 10, b)
}
func BenchmarkTcp_1_64k(b *testing.B) {
	runTestTcp(1, 16, b)
}
func BenchmarkTcp_2_64k(b *testing.B) {
	runTestTcp(2, 16, b)
}
func BenchmarkTcp_8_64k(b *testing.B) {
	runTestTcp(8, 16, b)
}
func BenchmarkTcp_16_64k(b *testing.B) {
	runTestTcp(16, 16, b)
}
func BenchmarkTcp_1_256k(b *testing.B) {
	runTestTcp(1, 18, b)
}
func BenchmarkTcp_2_256k(b *testing.B) {
	runTestTcp(2, 18, b)
}
func BenchmarkTcp_4_256k(b *testing.B) {
	runTestTcp(4, 18, b)
}
func BenchmarkTcp_8_256k(b *testing.B) {
	runTestTcp(8, 18, b)
}

func BenchmarkTcp_16_256k(b *testing.B) {
	runTestTcp(16, 18, b)
}
func BenchmarkTcp_1_512k(b *testing.B) {
	runTestTcp(1, 19, b)
}
func BenchmarkTcp_2_512k(b *testing.B) {
	runTestTcp(2, 19, b)
}
func BenchmarkTcp_16_512k(b *testing.B) {
	runTestTcp(16, 19, b)
}
func BenchmarkTcp_1_1m(b *testing.B) {
	runTestTcp(1, 20, b)
}
func BenchmarkTcp_2_1m(b *testing.B) {
	runTestTcp(2, 20, b)
}
