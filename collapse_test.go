package ipx_test

import (
	"fmt"
	"github.com/jwilner/ipx"
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
			[]string{
				"192.0.2.0/26",
				"192.0.2.64/26",
				"192.0.2.128/26",
				"192.0.2.192/26",
			},
			[]string{
				"192.0.2.0/24",
			},
		},
		{
			"simple v6",
			[]string{
				"::/26",
				"0:40::/26",
				"0:80::/26",
				"0:c0::/26",
			},
			[]string{
				"::/24",
			},
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
