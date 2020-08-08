# ipx

[![Tests](https://github.com/jwilner/ipx/workflows/tests/badge.svg)](https://github.com/jwilner/ipx/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwilner/ipx)](https://goreportcard.com/report/github.com/jwilner/ipx)
[![GoDoc](https://godoc.org/github.com/jwilner/ipx?status.svg)](https://godoc.org/github.com/jwilner/ipx)

Extending ip operations for Go to support common operations.

```go
func ExampleCollapse() {
	fmt.Println(ipx.Collapse(
		[]*net.IPNet{
			cidr("192.0.2.0/26"),
			cidr("192.0.2.64/26"),
			cidr("192.0.2.128/26"),
			cidr("192.0.2.192/26"),
		},
	))
	// Output:
	// [192.0.2.0/24]
}

func ExampleSplit() {
	c := cidr("10.0.0.0/24")
	splitter := ipx.Split(c, 26)
	for splitter.Next(c) {
		fmt.Println(c)
	}
	// Output:
	// 10.0.0.0/26
	// 10.0.0.64/26
	// 10.0.0.128/26
	// 10.0.0.192/26
}
```

See example tests for more usage.

## design thoughts

- Coordinate on stdlib types
- Avoid allocations whenever possible.
- Look to python ipaddress lib for feature list
