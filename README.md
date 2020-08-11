# ipx

[![Tests](https://github.com/ns1/ipx/workflows/tests/badge.svg)](https://github.com/ns1/ipx/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ns1/ipx)](https://goreportcard.com/report/github.com/ns1/ipx)
[![GoDoc](https://godoc.org/github.com/ns1/ipx?status.svg)](https://godoc.org/github.com/ns1/ipx)

Extending IP address support for Go.

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
	for splitter.Next() {
		fmt.Println(splitter.Net())
	}
	// Output:
	// 10.0.0.0/26
	// 10.0.0.64/26
	// 10.0.0.128/26
	// 10.0.0.192/26
}

func ExampleSummarizeRange() {
	fmt.Println(ipx.SummarizeRange(net.ParseIP("192.0.2.0"), net.ParseIP("192.0.2.130")))
	// Output:
	// [192.0.2.0/25 192.0.2.128/31 192.0.2.130/32]
}

func ExampleExclude() {
	fmt.Println(
		ipx.Exclude(
			cidr("10.1.1.0/24"),
			cidr("10.1.1.0/26"),
		),
	)
	// Output:
	// [10.1.1.128/25 10.1.1.64/26]
}
```

See example tests for more usage.

## design thoughts

- Coordinate on stdlib types
- Avoid allocations whenever possible.
- Look to python ipaddress lib for feature list
