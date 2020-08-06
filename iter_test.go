package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func ExampleIterIP() {
	ip := net.ParseIP("10.0.0.0")
	for i, iter := 0, ipx.IterIP(ip, 100); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(ip, -100); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
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
	for i, iter := 0, ipx.IterIP(ip, 1e18); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(ip, -1e18); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
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
	benchmarkIterIP(b, net.ParseIP("10.0.0.0"), 100)
}

func BenchmarkIterIP_IP6(b *testing.B) {
	benchmarkIterIP(b, net.ParseIP("::"), 100)
}

func benchmarkIterIP(b *testing.B, ip net.IP, step int) {
	b.ReportAllocs()

	tgt := make(net.IP, len(ip))

	b.ResetTimer()

	iter := ipx.IterIP(ip, step)
	for i := 0; i < b.N; i++ {
		if !iter.Next(tgt) {
			iter = ipx.IterIP(ip, step)
		}
	}
}

func ExampleIterNet() {
	ipN := cidr("10.0.0.0/16")
	for i, iter := 0, ipx.IterNet(ipN, 100); i < 5 && iter.Next(ipN); i++ {
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -100); i < 5 && iter.Next(ipN); i++ {
		fmt.Println(ipN)
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
	for i, iter := 0, ipx.IterNet(ipN, 1e18); i < 5 && iter.Next(ipN); i++ {
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -1e18); i < 5 && iter.Next(ipN); i++ {
		fmt.Println(ipN)
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
	benchmarkIterNet(b, cidr("10.0.0.0/16"), 100)
}

func BenchmarkIterNet_IP6(b *testing.B) {
	benchmarkIterNet(b, cidr("::/64"), 100)
}

func benchmarkIterNet(b *testing.B, ipN *net.IPNet, incr int) {
	b.ReportAllocs()

	n := &net.IPNet{IP: make(net.IP, len(ipN.IP)), Mask: make(net.IPMask, len(ipN.Mask))}

	b.ResetTimer()

	iter := ipx.IterNet(ipN, incr)
	for i := 0; i < b.N; i++ {
		if !iter.Next(n) {
			iter = ipx.IterNet(ipN, incr)
		}
	}
}
