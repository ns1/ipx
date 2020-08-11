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

func TestIncrIP(t *testing.T) {
	for _, c := range []struct {
		name, in string
		incr     int
		out      string
	}{
		{"ipv4 add", "0.0.0.0", 257, "0.0.1.1"},
		{"ipv4 minus", "0.0.1.1", -257, "0.0.0.0"},
		{"ipv6 add", "::", 1 << 32, "::1:0:0"},
		{"ipv6 minus", "::1:0:0", -(1 << 32), "::"},
	} {
		t.Run(c.name, func(t *testing.T) {
			ip := net.ParseIP(c.in)
			ipx.IncrIP(ip, c.incr)
			out := ip.String()
			if c.out != out {
				t.Fatalf("wanted %v but got %v", c.out, out)
			}
		})
	}

	t.Run("nil illegal", func(t *testing.T) {
		p := grabPanic(func() {
			ipx.IncrIP(nil, 1)
		})
		if p == nil {
			t.Fatal("expected a panic")
		}
		err, ok := p.(error)
		if !ok {
			t.Fatalf("expected an error in panic but got %T", err)
		}
		if err.Error() != "IP cannot be nil" {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func grabPanic(f func()) (p interface{}) {
	defer func() {
		p = recover()
	}()
	f()
	return
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

func TestIncrNet(t *testing.T) {
	for _, c := range []struct {
		name, in string
		incr     int
		out      string
	}{
		{"ipv4 add", "10.0.0.0/16", 2, "10.2.0.0/16"},
		{"ipv4 minus", "10.2.0.0/16", -2, "10.0.0.0/16"},
		{"ipv6 add", "::/32", 2, "0:2::/32"},
		{"ipv6 minus", "0:2::/32", -2, "::/32"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, in, _ := net.ParseCIDR(c.in)
			ipx.IncrNet(in, c.incr)
			out := in.String()
			if c.out != out {
				t.Fatalf("wanted %v but got %v", c.out, out)
			}
		})
	}

	t.Run("nil IP panics", func(t *testing.T) {
		p := grabPanic(func() {
			ipx.IncrNet(&net.IPNet{Mask: net.CIDRMask(32, 32)}, 1)
		})
		if p == nil {
			t.Fatal("expected panic")
		}
		err, ok := p.(error)
		if !ok {
			t.Fatalf("expected error but got %T", err)
		}
		if err.Error() != "IP cannot be nil" {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("nil mask panics", func(t *testing.T) {
		p := grabPanic(func() {
			ipx.IncrNet(&net.IPNet{IP: make(net.IP, 16)}, 1)
		})
		if p == nil {
			t.Fatal("expected panic")
		}
		err, ok := p.(error)
		if !ok {
			t.Fatalf("expected error but got %T", err)
		}
		if err.Error() != "mask cannot be nil" {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
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
