package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func BenchmarkIncrIP_IP4(b *testing.B) {
	b.ReportAllocs()

	ip := net.ParseIP("0.0.0.0")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ipx.IncrIP(ip, 1)
	}
}

func ExampleIncrIP() {
	ip := net.ParseIP("0.0.0.0")
	ipx.IncrIP(ip, 257)
	fmt.Println(ip)
	// Output:
	// 0.0.1.1
}

func BenchmarkIncrIP_IP6(b *testing.B) {
	b.ReportAllocs()

	ip := net.ParseIP("::")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ipx.IncrIP(ip, 1)
	}
}

func ExampleIncrIP_IP6() {
	ip := net.ParseIP("::")
	ipx.IncrIP(ip, 1<<32)
	fmt.Println(ip)
	// Output:
	// ::1:0:0
}

func BenchmarkIncrNet_IP4(b *testing.B) {
	b.ReportAllocs()

	ipN := cidr("10.0.0.0/16")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ipx.IncrNet(ipN, 1)
	}
}

func ExampleIncrNet() {
	ipN := cidr("10.0.0.0/16")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN)
	// Output:
	// 10.2.0.0/16
}

func BenchmarkIncrNet_IP6(b *testing.B) {
	b.ReportAllocs()

	ipN := cidr("::/32")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ipx.IncrNet(ipN, 1)
	}
}

func ExampleIncrNet_IP6() {
	ipN := cidr("::/32")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN)
	// Output:
	// 0:2::/32
}

func cidr(cidrS string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrS)
	return ipNet
}
