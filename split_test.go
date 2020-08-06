package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
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
