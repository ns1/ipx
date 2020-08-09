package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"net"
	"testing"
)

func ExampleSummarizeRange() {
	fmt.Println(ipx.SummarizeRange(net.ParseIP("192.0.2.0"), net.ParseIP("192.0.2.130")))
	// Output:
	// [192.0.2.0/25 192.0.2.128/31 192.0.2.130/32]
}

func ExampleSummarizeRange_IP6() {
	fmt.Println(ipx.SummarizeRange(
		net.ParseIP("::"),
		net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
	))
	// Output:
	// [::/0]
}

func BenchmarkSummarizeRange(b *testing.B) {
	type bench struct {
		name        string
		first, last string
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{
					"all",
					"0.0.0.0",
					"255.255.255.255",
				},
				{
					"32",
					"0.0.0.1",
					"255.255.255.255",
				},
			},
		},
		{
			"ipv6",
			[]bench{
				{
					"all",
					"::",
					"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
				},
				{
					"128",
					"::1",
					"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
				},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				first, last := net.ParseIP(c.first), net.ParseIP(c.last)
				b.Run(c.name, func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = ipx.SummarizeRange(first, last)
					}
				})
			}
		})
	}
}

func TestSummarizeRange(t *testing.T) {
	for _, c := range []struct {
		name        string
		first, last string
		results     []string
	}{
		{"no overlap", "192.0.2.1", "192.0.2.0", []string{}},
		{
			"simple",
			"192.0.2.0",
			"192.0.2.130",
			[]string{"192.0.2.0/25", "192.0.2.128/31", "192.0.2.130/32"},
		},
		{
			"single",
			"192.0.2.0",
			"192.0.2.0",
			[]string{"192.0.2.0/32"},
		},
		{
			"all",
			"0.0.0.0",
			"255.255.255.255",
			[]string{"0.0.0.0/0"},
		},
		{
			"odd start",
			"192.0.2.101",
			"192.0.2.130",
			[]string{
				"192.0.2.101/32",
				"192.0.2.102/31",
				"192.0.2.104/29",
				"192.0.2.112/28",
				"192.0.2.128/31",
				"192.0.2.130/32",
			},
		},
		{
			"ipv6",
			"1::",
			"1:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			[]string{"1::/16"},
		},
		{
			"ipv6 all",
			"::",
			"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			[]string{"::/0"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			var found []string
			for _, ipN := range ipx.SummarizeRange(net.ParseIP(c.first), net.ParseIP(c.last)) {
				found = append(found, ipN.String())
			}
			if len(c.results) != len(found) {
				t.Fatalf("expected %v elements but got %v: %v", len(c.results), len(found), found)
			}
			for i := range found {
				if c.results[i] != found[i] {
					t.Errorf("position %d: expected %v but got %v", i, c.results[i], found[i])
				}
			}
		})
	}
}
