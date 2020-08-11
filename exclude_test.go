package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
	"testing"
)

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

func TestExclude(t *testing.T) {
	for _, c := range []struct {
		name, a, b string
		expected   []string
	}{
		{"disjoint", "10.1.1.0/24", "10.0.1.0/26", []string{"10.1.1.0/24"}},
		{"mismatch versions", "10.1.1.0/24", "2001:db8::1/128", []string{"10.1.1.0/24"}},
		{"ipv4", "10.1.1.0/24", "10.1.1.0/26", []string{"10.1.1.128/25", "10.1.1.64/26"}},
		{
			"ipv4 leaf", "10.1.1.0/24", "10.1.1.5/32",
			[]string{
				"10.1.1.128/25",
				"10.1.1.64/26",
				"10.1.1.32/27",
				"10.1.1.16/28",
				"10.1.1.8/29",
				"10.1.1.0/30",
				"10.1.1.6/31",
				"10.1.1.4/32",
			},
		},
		{
			"ipv6", "2001:db8::/124", "2001:db8::1/128",
			[]string{
				"2001:db8::8/125",
				"2001:db8::4/126",
				"2001:db8::2/127",
				"2001:db8::/128",
			},
		},
		{"ipv6 left", "2001:db8::/124", "2001:db8::8/126", []string{"2001:db8::/125", "2001:db8::c/126"}},
		{"ipv6 right", "2001:db8::/124", "2001:db8::0/126", []string{"2001:db8::8/125", "2001:db8::4/126"}},
	} {
		t.Run(c.name, func(t *testing.T) {
			var got []string
			for _, n := range ipx.Exclude(cidr(c.a), cidr(c.b)) {
				got = append(got, n.String())
			}
			if len(c.expected) != len(got) {
				t.Fatalf("wanted %v but got %v items: %v", len(c.expected), len(got), got)
			}
			for i := range got {
				if c.expected[i] != got[i] {
					t.Errorf("wanted %v at position %d but got %v", c.expected[i], i, got[i])
				}
			}
		})
	}
}

func BenchmarkExclude(b *testing.B) {
	type bench struct {
		a, b string
	}
	for _, g := range []struct {
		name    string
		benches []bench
	}{
		{
			"ipv4",
			[]bench{
				{"10.1.1.0/24", "10.1.1.0/26"},
				{"0.0.0.0/0", "10.1.1.0/32"},
			},
		},
		{
			"ipv6",
			[]bench{
				{"2001:db8::1/124", "2001:db8::1/126"},
				{"::/0", "2001:db8::1/128"},
			},
		},
	} {
		b.Run(g.name, func(b *testing.B) {
			for _, c := range g.benches {
				netA, netB := cidr(c.a), cidr(c.b)
				onesA, _ := netA.Mask.Size()
				onesB, _ := netB.Mask.Size()
				b.Run(fmt.Sprintf("%v-%v", onesA, onesB), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = ipx.Exclude(netA, netB)
					}
				})
			}
		})
	}
}
