package ipx

import (
	"net"
)

const (
	maxUint32 = 1<<32 - 1
	maxUint64 = 1<<64 - 1
)

// IPIter permits iteration over a series of ips. It is always start inclusive.
type IPIter interface {
	// Next returns true when the underlying pointer has been successfully updated
	// with the next value.
	Next(ip net.IP) bool
}

// NetIter permits iteration over a series of IP networks. It is always start inclusive.
type NetIter interface {
	// Next returns true when the underlying pointer has been successfully updated
	// with the next value.
	Next(ipNet *net.IPNet) bool
}

// IterIP returns an Iter for the given increment over the IPs
func IterIP(ip net.IP, incr int) IPIter {
	is4 := ip.To4() != nil
	sub := incr < 0
	switch {
	case is4 && sub:
		return &decrIP4{to32(ip), uint32(incr * -1), 0}
	case is4:
		return &incrIP4{to32(ip), uint32(incr), maxUint32}
	case sub:
		return &decrIP6{to128(ip), uint128{0, uint64(incr * -1)}, uint128{}}
	default:
		return &incrIP6{
			to128(ip),
			uint128{0, uint64(incr)},
			uint128{maxUint64, maxUint64},
		}
	}
}

// IterNet returns an iterator for the given increment starting with the provided network
func IterNet(ipNet *net.IPNet, incr int) NetIter {
	ones, bits := ipNet.Mask.Size()
	suffix := bits - ones

	sub := incr < 0
	is4 := ipNet.IP.To4() != nil

	switch {
	case is4 && sub:
		return &decrIP4Net{
			decrIP4{
				to32(ipNet.IP),
				uint32(incr*-1) << suffix,
				0,
			},
			to32(ipNet.Mask),
		}
	case is4:
		return &incrIP4Net{
			incrIP4{
				to32(ipNet.IP),
				uint32(incr) << suffix,
				maxUint32,
			},
			to32(ipNet.Mask),
		}
	case sub:
		return &decrIP6Net{
			decrIP6{
				to128(ipNet.IP),
				uint128{0, uint64(incr * -1)}.Lsh(uint(suffix)),
				uint128{},
			},
			to128(ipNet.Mask),
		}
	default:
		return &incrIP6Net{
			incrIP6{
				to128(ipNet.IP),
				uint128{0, uint64(incr)}.Lsh(uint(suffix)),
				uint128{maxUint64, maxUint64},
			},
			to128(ipNet.Mask),
		}
	}
}

type decrIP4 struct {
	v, decr, limit uint32
}

func (d *decrIP4) Next(ip net.IP) bool {
	if d.v <= d.limit {
		return false
	}
	from32(d.v, ip)
	d.v -= d.decr
	return true
}

type decrIP4Net struct {
	decrIP4
	mask uint32
}

func (d *decrIP4Net) Next(ipN *net.IPNet) bool {
	if !d.decrIP4.Next(ipN.IP) {
		return false
	}
	from32(d.mask, ipN.Mask)
	return true
}

type incrIP4 struct {
	v, incr, limit uint32
}

func (i *incrIP4) Next(ip net.IP) bool {
	if i.v >= i.limit {
		return false
	}
	from32(i.v, ip)
	i.v += i.incr
	return true
}

type incrIP4Net struct {
	incrIP4
	mask uint32
}

func (i *incrIP4Net) Next(ipN *net.IPNet) bool {
	if !i.incrIP4.Next(ipN.IP) {
		return false
	}
	from32(i.mask, ipN.Mask)
	return true
}

type incrIP6 struct {
	v, incr, limit uint128
}

func (i *incrIP6) Next(ip net.IP) bool {
	if i.v.Cmp(i.limit) != -1 {
		return false
	}
	from128(i.v, ip)
	i.v = i.v.Add(i.incr)
	return true
}

type incrIP6Net struct {
	incrIP6
	mask uint128
}

func (i *incrIP6Net) Next(ipN *net.IPNet) bool {
	if !i.incrIP6.Next(ipN.IP) {
		return false
	}
	from128(i.mask, ipN.Mask)
	return true
}

type decrIP6 struct {
	v, decr, limit uint128
}

func (d *decrIP6) Next(ip net.IP) bool {
	if d.v.Cmp(d.limit) == -1 {
		return false
	}
	from128(d.v, ip)
	d.v = d.v.Minus(d.decr)
	return true
}

type decrIP6Net struct {
	decrIP6
	mask uint128
}

func (d *decrIP6Net) Next(ipN *net.IPNet) bool {
	if !d.decrIP6.Next(ipN.IP) {
		return false
	}
	from128(d.mask, ipN.Mask)
	return true
}
