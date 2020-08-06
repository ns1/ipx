# ipx

Extending ip operations for Go to support common operations.

```go
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
