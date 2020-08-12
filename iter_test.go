package ipx_test

import (
	"fmt"
	"github.com/ns1/ipx"
	"net"
	"testing"
)

func TestIterIP(t *testing.T) {
	for _, c := range []struct {
		name     string
		start    string
		step     int
		end      string
		expected []string
	}{
		{"ipv4 incr", "10.0.0.0", 1, "10.0.0.3", []string{"10.0.0.0", "10.0.0.1", "10.0.0.2"}},
		{"ipv4 incr equal", "10.0.0.0", 1, "10.0.0.0", []string{}},
		{"ipv4 incr end exclusive", "10.0.0.0", 2, "10.0.0.2", []string{"10.0.0.0"}},
		{"ipv4 incr end awkward", "10.0.0.0", 2, "10.0.0.3", []string{"10.0.0.0", "10.0.0.2"}},

		{"ipv4 decr", "10.0.0.3", -1, "10.0.0.0", []string{"10.0.0.3", "10.0.0.2", "10.0.0.1"}},
		{"ipv4 decr equal", "10.0.0.0", -1, "10.0.0.0", []string{}},
		{"ipv4 decr end exclusive", "10.0.0.3", -2, "10.0.0.1", []string{"10.0.0.3"}},
		{"ipv4 decr end awkward", "10.0.0.3", -2, "10.0.0.0", []string{"10.0.0.3", "10.0.0.1"}},

		{"ipv6 incr", "::", 1, "::3", []string{"::", "::1", "::2"}},
		{"ipv6 incr equal", "::", 1, "::", []string{}},
		{"ipv6 incr end exclusive", "::", 2, "::2", []string{"::"}},
		{"ipv6 incr end awkward", "::", 2, "::3", []string{"::", "::2"}},

		{"ipv6 decr", "::3", -1, "::", []string{"::3", "::2", "::1"}},
		{"ipv6 decr equal", "::", -1, "::", []string{}},
		{"ipv6 decr end exclusive", "::3", -2, "::1", []string{"::3"}},
		{"ipv6 decr end awkward", "::3", -2, "::", []string{"::3", "::1"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			//if !strings.HasSuffix(t.Name(), "ipv4_decr_end_awkward") {
			//	t.SkipNow()
			//}
			start, end := net.ParseIP(c.start), net.ParseIP(c.end)
			var got []string
			for ips := ipx.IterIP(start, c.step, end); ips.Next(); {
				got = append(got, ips.IP().String())
			}
			if len(got) != len(c.expected) {
				t.Fatalf("Expected %v items but got %v: %v", len(c.expected), len(got), got)
			}
			for i := range got {
				if c.expected[i] != got[i] {
					t.Errorf("expected %v at position %v but got %v", c.expected[i], i, got[i])
				}
			}
		})
	}
}

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
	ip := net.ParseIP("2001:db8::")
	for i, iter := 0, ipx.IterIP(ip, 1e18, nil); i < 5 && iter.Next(); i++ {
		ip = iter.IP()
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(ip, -1e18, nil); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.IP())
	}
	// Output:
	// 2001:db8::
	// 2001:db8::de0:b6b3:a764:0
	// 2001:db8::1bc1:6d67:4ec8:0
	// 2001:db8::29a2:241a:f62c:0
	// 2001:db8::3782:dace:9d90:0
	// 2001:db8::3782:dace:9d90:0
	// 2001:db8::29a2:241a:f62c:0
	// 2001:db8::1bc1:6d67:4ec8:0
	// 2001:db8::de0:b6b3:a764:0
	// 2001:db8::
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
	for i, iter := 0, ipx.IterNet(ipN, 100, nil); i < 5 && iter.Next(); i++ {
		ipN = iter.Net()
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -100, nil); i < 5 && iter.Next(); i++ {
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
	ipN := cidr("2001:db8::/64")
	for i, iter := 0, ipx.IterNet(ipN, 1e18, nil); i < 5 && iter.Next(); i++ {
		ipN = iter.Net()
		fmt.Println(ipN)
	}
	for i, iter := 0, ipx.IterNet(ipN, -1e18, nil); i < 5 && iter.Next(); i++ {
		fmt.Println(iter.Net())
	}
	// Output:
	// 2001:db8::/64
	// 2de1:c46b:a764::/64
	// 3bc2:7b1f:4ec8::/64
	// 49a3:31d2:f62c::/64
	// 5783:e886:9d90::/64
	// 5783:e886:9d90::/64
	// 49a3:31d2:f62c::/64
	// 3bc2:7b1f:4ec8::/64
	// 2de1:c46b:a764::/64
	// 2001:db8::/64
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

					iter := ipx.IterNet(ipN, c.incr, nil)
					for i := 0; i < b.N; i++ {
						if !iter.Next() {
							iter = ipx.IterNet(ipN, c.incr, nil)
						}
					}
				})
			}
		})
	}
}
