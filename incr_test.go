package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
)

func ExampleIncrIP() {
	ip := net.ParseIP("0.0.0.0")
	ipx.IncrIP(ip, 257)
	fmt.Println(ip.String())
	// Output:
	// 0.0.1.1
}

func ExampleIncrIP_IP6() {
	ip := net.ParseIP("::")
	ipx.IncrIP(ip, 1<<32)
	fmt.Println(ip.String())
	// Output:
	// ::1:0:0
}

func ExampleIncrNet() {
	ipN := cidr("10.0.0.0/16")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN.String())
	// Output:
	// 10.2.0.0/16
}

func ExampleIncrNet_IP6() {
	ipN := cidr("::/32")
	ipx.IncrNet(ipN, 2)
	fmt.Println(ipN.String())
	// Output:
	// 0:2::/32
}

func cidr(cidrS string) net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidrS)
	return *ipNet
}
