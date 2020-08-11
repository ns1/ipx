package ipx_test

import (
	"fmt"
	"github.com/ns1/ipx"
	"net"
	"testing"
)

func ExampleIterIP() {
	ip := net.ParseIP("10.0.0.0")
	for i, iter := 0, ipx.IterIP(ip, 100, nil); i < 5 && iter.Next(); i++ {
		ip = iter.IP()
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(ip, -100, nil); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.IP())
	}
	// Output:
	// 10.0.0.0
	// 10.0.0.100
	// 10.0.0.200
	// 10.0.1.44
	// 10.0.1.144
	// 10.0.1.144
	// 10.0.1.44
	// 10.0.0.200
	// 10.0.0.100
	// 10.0.0.0
}

func ExampleIterIP_IP6() {
	ip := net.ParseIP("::")
	for i, iter := 0, ipx.IterIP(ip, 1e18, nil); i < 5 && iter.Next(); i++ {
		ip = iter.IP()
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(ip, -1e18, nil); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.IP())
	}
	// Output:
	// ::
	// ::de0:b6b3:a764:0
	// ::1bc1:6d67:4ec8:0
	// ::29a2:241a:f62c:0
	// ::3782:dace:9d90:0
	// ::3782:dace:9d90:0
	// ::29a2:241a:f62c:0
	// ::1bc1:6d67:4ec8:0
	// ::de0:b6b3:a764:0
	// ::
}

func BenchmarkIterIP(b *testing.B) {
	type bench struct {
		ip   string
		incr int
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{"10.0.0.0", 100},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::", 100},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				ip := net.ParseIP(c.ip)
				b.Run(fmt.Sprint(c.incr), func(b *testing.B) {
					b.ReportAllocs()

					iter := ipx.IterIP(ip, c.incr, nil)
					for i := 0; i < b.N; i++ {
						if !iter.Next() {
							iter = ipx.IterIP(iter.IP(), c.incr, nil)
						}
					}
				})
			}
		})
	}
}

func ExampleIterNet() {
	ipN := cidr("10.0.0.0/16")
	for i, iter := 0, ipx.IterNet(ipN, 100); i < 5 && iter.Next(); i++ {
		ipN = iter.Net()
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -100); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.Net())
	}
	// Output:
	// 10.0.0.0/16
	// 10.100.0.0/16
	// 10.200.0.0/16
	// 11.44.0.0/16
	// 11.144.0.0/16
	// 11.144.0.0/16
	// 11.44.0.0/16
	// 10.200.0.0/16
	// 10.100.0.0/16
	// 10.0.0.0/16
}

func ExampleIterNet_IP6() {
	ipN := cidr("::/64")
	for i, iter := 0, ipx.IterNet(ipN, 1e18); i < 5 && iter.Next(); i++ {
		ipN = iter.Net()
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -1e18); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.Net())
	}
	// Output:
	// ::/64
	// de0:b6b3:a764::/64
	// 1bc1:6d67:4ec8::/64
	// 29a2:241a:f62c::/64
	// 3782:dace:9d90::/64
	// 3782:dace:9d90::/64
	// 29a2:241a:f62c::/64
	// 1bc1:6d67:4ec8::/64
	// de0:b6b3:a764::/64
	// ::/64
}

func BenchmarkIterNet(b *testing.B) {
	type bench struct {
		cidr string
		incr int
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{"10.0.0.0/16", 100},
				{"192.0.2.0/24", 1},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::/64", 100},
				{"::/120", 1},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				ipN := cidr(c.cidr)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprintf("%v-%v", ones, c.incr), func(b *testing.B) {
					b.ReportAllocs()

					iter := ipx.IterNet(ipN, c.incr)
					for i := 0; i < b.N; i++ {
						if !iter.Next() {
							iter = ipx.IterNet(ipN, c.incr)
						}
					}
				})
			}
		})
	}
}
