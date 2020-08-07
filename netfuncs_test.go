package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func ExampleSupernet() {
	ipN := cidr("192.0.2.0/24")
	super := ipx.Supernet(ipN, 20)
	fmt.Println(super)
	// Output:
	// 192.0.0.0/20
}

func BenchmarkSupernet(b *testing.B) {
	type bench struct {
		cidr      string
		newPrefix int
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{"192.0.2.0/24", 20},
				{"192.0.2.0/24", 15},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::/24", 20},
				{"::/24", 15},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				ipN := cidr(c.cidr)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprintf("%v-%v", ones, c.newPrefix), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						_ = ipx.Supernet(ipN, c.newPrefix)
					}
				})
			}
		})
	}
}

func ExampleBroadcast() {
	ipN := cidr("10.0.1.0/24")
	fmt.Println(ipx.Broadcast(ipN))
	// Output:
	// 10.0.1.255
}

func BenchmarkBroadcast(b *testing.B) {
	for _, g := range []struct {
		name  string
		cidrs []*net.IPNet
	}{
		{
			"ipv4",
			[]*net.IPNet{
				cidr("10.0.1.0/31"),
				cidr("10.1.0.0/16"),
				cidr("0.0.0.0/0"),
			},
		},
		{
			"ipv6",
			[]*net.IPNet{
				cidr("::/127"),
				cidr("::/64"),
				cidr("::/0"),
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cidrs {
				ones, _ := c.Mask.Size()
				b.Run(fmt.Sprint(ones), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = ipx.Broadcast(c)
					}
				})
			}
		})
	}
}

func ExampleIsSubnet() {
	a, b := cidr("10.0.0.0/16"), cidr("10.0.1.0/24")
	fmt.Println(ipx.IsSubnet(a, b))
	fmt.Println(ipx.IsSubnet(a, a))
	fmt.Println(ipx.IsSubnet(b, a))
	// Output:
	// true
	// true
	// false
}

func ExampleIsSupernet() {
	a, b := cidr("10.0.1.0/24"), cidr("10.0.0.0/16")
	fmt.Println(ipx.IsSupernet(a, b))
	fmt.Println(ipx.IsSupernet(a, a))
	fmt.Println(ipx.IsSupernet(b, a))
	// Output:
	// true
	// true
	// false
}
