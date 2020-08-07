package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func BenchmarkIncrIP(b *testing.B) {
	type bench struct {
		ip   string
		incr int
	}
	for _, g := range []struct {
		name  string
		cases []bench
	}{
		{
			"ipv4",
			[]bench{
				{"10.0.0.0", 1},
				{"10.0.0.0", 2},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::", 1},
				{"::", 2},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cases {
				ip := net.ParseIP(c.ip)
				b.Run(fmt.Sprint(c.incr), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						ipx.IncrIP(ip, c.incr)
					}
				})
			}
		})
	}

}

func ExampleIncrIP() {
	ip := net.ParseIP("0.0.0.0")
	ipx.IncrIP(ip, 257)
	fmt.Println(ip)
	// Output:
	// 0.0.1.1
}

func ExampleIncrIP_IP6() {
	ip := net.ParseIP("::")
	ipx.IncrIP(ip, 1<<32)
	fmt.Println(ip)
	// Output:
	// ::1:0:0
}

func ExampleIncrNet() {
	ipN := cidr("10.0.0.0/16")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN)
	// Output:
	// 10.2.0.0/16
}

func ExampleIncrNet_IP6() {
	ipN := cidr("::/32")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN)
	// Output:
	// 0:2::/32
}

func BenchmarkIncrNet(b *testing.B) {
	type bench struct {
		cidr string
		incr int
	}
	for _, g := range []struct {
		name  string
		cases []bench
	}{
		{
			"ipv4",
			[]bench{
				{"10.0.0.0/30", 1},
				{"10.0.0.0/30", 2},
				{"10.0.0.0/24", 1},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::/126", 1},
				{"::/32", 1},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cases {
				ipN := cidr(c.cidr)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprintf("%v-%v", ones, c.incr), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						ipx.IncrNet(ipN, c.incr)
					}
				})
			}
		})
	}
}

func cidr(cidrS string) *net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrS)
	return ipNet
}
