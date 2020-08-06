package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
)

func ExampleSupernet() {
	ipN := cidr("192.0.2.0/24")
	super := ipx.Supernet(ipN, 20)
	fmt.Println(super)
	// Output:
	// 192.0.0.0/20
}

func ExampleBroadcast() {
	ipN := cidr("10.0.1.0/24")
	fmt.Println(ipx.Broadcast(ipN))
	// Output:
	// 10.0.1.255
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
