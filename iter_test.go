package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
)

func ExampleIterIP() {
	ip := net.ParseIP("10.0.0.0")
	for i, iter := 0, ipx.IterIP(100); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(-100); i < 5 && iter.Next(ip); i++ {
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
	for i, iter := 0, ipx.IterIP(1e18); i < 5 && iter.Next(ip); i++ {
		fmt.Println(ip)
	}
	for i, iter := 0, ipx.IterIP(-1e18); i < 5 && iter.Next(ip); i++ {
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
