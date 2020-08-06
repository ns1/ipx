package ipx

import (
	"math"
	"net"
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
	return iterIP(ip, incr)
}

func iterIP(ip net.IP, incr int) IPIter {
	is4 := ip.To4() != nil
	sub := incr < 0
	switch {
	case is4 && sub:
		return &decrIP4{to32(ip), uint32(incr * -1), 0}
	case is4:
		return &incrIP4{to32(ip), uint32(incr), math.MaxUint32}
	case sub:
		return &decrIP6{toUint128(ip), uint128{0, uint64(incr * -1)}, uint128{}}
	default:
		return &incrIP6{
			toUint128(ip),
			uint128{0, uint64(incr)},
			uint128{math.MaxUint64, math.MaxUint64},
		}
	}
}

// IterNet returns an iterator for the given increment starting with the provided network
func IterNet(ipNet *net.IPNet, incr int) NetIter {
	return iterNet(ipNet.IP, ipNet.Mask, incr)
}

func iterNet(ip net.IP, mask net.IPMask, incr int) NetIter {
	ipIter := func() IPIter {
		ones, bits := mask.Size()
		suffix := bits - ones

		sub := incr < 0
		is4 := ip.To4() != nil

		switch {
		case is4 && sub:
			return &decrIP4{to32(ip), uint32(incr*-1) << suffix, 0}
		case is4:
			return &incrIP4{to32(ip), uint32(incr) << suffix, math.MaxUint32}
		case sub:
			decrB := uint128{0, uint64(incr * -1)}
			decrB.Lsh(uint(suffix))
			return &decrIP6{toUint128(ip), decrB, uint128{}}
		default:
			incrB := uint128{0, uint64(incr)}
			incrB.Lsh(uint(suffix))
			return &incrIP6{toUint128(ip), incrB, uint128{math.MaxUint64, math.MaxUint64}}
		}
	}()
	return &netIter{IPIter: ipIter, mask: mask}
}

type netIter struct {
	IPIter
	mask net.IPMask
}

func (n *netIter) Next(ipNet *net.IPNet) bool {
	if !n.IPIter.Next(ipNet.IP) {
		return false
	}
	copy(ipNet.Mask, n.mask)
	return true
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

type incrIP6 struct {
	v, incr, limit uint128
}

func (i *incrIP6) Next(ip net.IP) bool {
	if i.v.Cmp(i.limit) != -1 {
		return false
	}
	fromUint128(i.v, ip)
	i.v.Add(i.incr)
	return true
}

type decrIP6 struct {
	v, decr, limit uint128
}

func (d *decrIP6) Next(ip net.IP) bool {
	if d.v.Cmp(d.limit) == -1 {
		return false
	}
	fromUint128(d.v, ip)
	d.v.Minus(d.decr)
	return true
}
