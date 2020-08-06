package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
)

func ExampleSplit() {
	c := cidr("10.0.0.0/24")
	split := ipx.Split(c, 26)
	for split.Next(c) {
		fmt.Println(c.String())
	}
	// Output:
	// 10.0.0.0/26
	// 10.0.0.64/26
	// 10.0.0.128/26
	// 10.0.0.192/26
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
