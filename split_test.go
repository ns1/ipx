package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func ExampleSplit() {
	c := cidr("10.0.0.0/24")
	split := ipx.Split(c, 26)
	for split.Next(c) {
		fmt.Println(c)
	}
	// Output:
	// 10.0.0.0/26
	// 10.0.0.64/26
	// 10.0.0.128/26
	// 10.0.0.192/26
}

func ExampleSplit_IP6() {
	c := cidr("::/24")
	split := ipx.Split(c, 26)
	for split.Next(c) {
		fmt.Println(c)
	}
	// Output:
	// ::/26
	// 0:40::/26
	// 0:80::/26
	// 0:c0::/26
}

func BenchmarkSplit(b *testing.B) {
	b.ReportAllocs()

	c := cidr("10.0.0.0/24")

	ipN := net.IPNet{IP: make(net.IP, len(c.IP)), Mask: make(net.IPMask, len(c.Mask))}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for split := ipx.Split(c, 26); split.Next(&ipN); {
		}
	}
}

func ExampleAddresses() {
	c := cidr("10.0.0.0/30")
	addrs := ipx.Addresses(c)
	ip := make(net.IP, net.IPv4len)
	for addrs.Next(ip) {
		fmt.Println(ip)
	}
	// Output:
	// 10.0.0.0
	// 10.0.0.1
	// 10.0.0.2
	// 10.0.0.3
}

func BenchmarkAddresses(b *testing.B) {
	b.ReportAllocs()

	c := cidr("10.0.0.0/30")

	ip := make(net.IP, len(c.IP))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for addrs := ipx.Addresses(c); addrs.Next(ip); {
		}
	}
}

func ExampleHosts() {
	c := cidr("10.0.0.0/29")
	hosts := ipx.Hosts(c)
	ip := make(net.IP, net.IPv4len)
	for hosts.Next(ip) {
		fmt.Println(ip)
	}
	// Output:
	// 10.0.0.1
	// 10.0.0.2
	// 10.0.0.3
	// 10.0.0.4
	// 10.0.0.5
	// 10.0.0.6
}

func ExampleHosts_IP6() {
	c := cidr("::/125")
	hosts := ipx.Hosts(c)
	ip := make(net.IP, len(c.IP))
	for hosts.Next(ip) {
		fmt.Println(ip)
	}
	// Output:
	// ::1
	// ::2
	// ::3
	// ::4
	// ::5
	// ::6
}

func BenchmarkHosts(b *testing.B) {
	b.ReportAllocs()

	c := cidr("10.0.0.0/30")

	ip := make(net.IP, len(c.IP))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for hosts := ipx.Hosts(c); hosts.Next(ip); {
		}
	}
}
