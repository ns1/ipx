package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"testing"
)

func ExampleSplit() {
	c := cidr("10.0.0.0/24")
	split := ipx.Split(c, 26)
	for split.Next() {
		fmt.Println(split.Net())
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
	for split.Next() {
		fmt.Println(split.Net())
	}
	// Output:
	// ::/26
	// 0:40::/26
	// 0:80::/26
	// 0:c0::/26
}

func BenchmarkSplit(b *testing.B) {
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
				{"192.0.2.0/24", 30},
				{"192.0.2.0/24", 28},
			},
		},
		{
			"ipv6",
			[]bench{
				{"::/24", 30},
				{"::/24", 28},
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
						for split := ipx.Split(ipN, c.newPrefix); split.Next(); {
						}
					}

				})
			}
		})
	}
}

func ExampleAddresses() {
	c := cidr("10.0.0.0/30")
	addrs := ipx.Addresses(c)
	for addrs.Next() {
		fmt.Println(addrs.IP())
	}
	// Output:
	// 10.0.0.0
	// 10.0.0.1
	// 10.0.0.2
	// 10.0.0.3
}

func BenchmarkAddresses(b *testing.B) {
	for _, g := range []struct {
		name  string
		cidrs []string
	}{
		{
			"ipv4",
			[]string{
				"10.0.0.0/30", // 4
				"10.0.0.0/24", // 256
			},
		},
		{
			"ipv6",
			[]string{
				"::/126", // 4
				"::/120", // 256
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cidrs {
				ipN := cidr(c)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprint(ones), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						for hosts := ipx.Addresses(ipN); hosts.Next(); {
						}
					}
				})
			}
		})
	}
}

func ExampleHosts() {
	c := cidr("10.0.0.0/29")
	hosts := ipx.Hosts(c)
	for hosts.Next() {
		fmt.Println(hosts.IP())
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
	for hosts.Next() {
		fmt.Println(hosts.IP())
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
	for _, g := range []struct {
		name  string
		cidrs []string
	}{
		{
			"ipv4",
			[]string{
				"10.0.0.0/30", // 4-2
				"10.0.0.0/24", // 256-2
			},
		},
		{
			"ipv6",
			[]string{
				"::/126", // 4-2
				"::/120", // 256-2
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.cidrs {
				ipN := cidr(c)
				ones, _ := ipN.Mask.Size()
				b.Run(fmt.Sprint(ones), func(b *testing.B) {
					b.ReportAllocs()

					for i := 0; i < b.N; i++ {
						for hosts := ipx.Hosts(ipN); hosts.Next(); {
						}
					}
				})
			}
		})
	}
}
