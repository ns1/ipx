package ipx_test

import (
	"fmt"
	"github.com/ns1/ipx"
	"net"
	"testing"
)

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

func TestCollapse(t *testing.T) {
	for _, c := range []struct {
		name string
		in   []string
		out  []string
	}{
		{"empty", nil, nil},
		{
			"simple",
			[]string{"192.0.2.0/26", "192.0.2.64/26", "192.0.2.128/26", "192.0.2.192/26"},
			[]string{"192.0.2.0/24"},
		},
		{
			"dupe",
			[]string{"192.0.2.0/26", "192.0.2.64/26", "192.0.2.128/26", "192.0.2.192/26", "192.0.2.192/26"},
			[]string{"192.0.2.0/24"},
		},
		{
			"simple v6",
			[]string{"::/26", "0:40::/26", "0:80::/26", "0:c0::/26"},
			[]string{"::/24"},
		},
		{
			"dupe v6",
			[]string{"::/26", "0:40::/26", "0:80::/26", "0:c0::/26", "0:c0::/26"},
			[]string{"::/24"},
		},
		{
			"multi type",
			[]string{"0:80::/26", "0:c0::/26", "192.0.2.0/26", "192.0.2.64/26"},
			[]string{"192.0.2.0/25", "0:80::/25"},
		},
		{
			"ipv4 child included",
			[]string{"192.0.2.0/26", "192.0.2.64/26", "192.0.2.64/27"},
			[]string{"192.0.2.0/25"},
		},
		{
			"ipv4 disjoint",
			[]string{"192.0.2.0/27", "192.0.2.64/27", "192.0.2.64/27"},
			[]string{"192.0.2.0/27", "192.0.2.64/27"},
		},
		{
			"ipv6 child included",
			[]string{"0:80::/26", "0:c0::/26", "0:c0::/27"},
			[]string{"0:80::/25"},
		},
		{
			"ipv6 disjoint",
			[]string{"0:80::/27", "0:c0::/27", "0:c0::/27"},
			[]string{"0:80::/27", "0:c0::/27"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			in := make([]*net.IPNet, 0, len(c.in))
			for _, s := range c.in {
				in = append(in, cidr(s))
			}
			out := ipx.Collapse(in)
			got := make([]string, 0, len(out))
			for _, o := range out {
				got = append(got, o.String())
			}

			if len(c.out) != len(got) {
				t.Fatalf("Wanted length %v but got %v: %v", len(c.out), len(got), got)
			}

			for i := range c.out {
				if c.out[i] != got[i] {
					t.Errorf("Wanted %v but got %v at position %v", c.out[i], got[i], i)
				}
			}
		})
	}
}
