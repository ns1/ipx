package ipx

import (
	"net"
)

const (
	maxUint32 = 1<<32 - 1
	maxUint64 = 1<<64 - 1
)

type v4IPIter struct {
	val, incr, limit uint32
}

type v6IPIter struct {
	val, incr, limit uint128
}

const (
	ipIterFlagV6 = 1 << iota
	ipIterFlagNegative
)

// IPIter permits iteration over a series of ips. It is always start inclusive.
type IPIter struct {
	v4    v4IPIter
	v6    v6IPIter
	flags uint8
}

// Next returns true when the underlying pointer has been successfully updated with the next value.
func (i *IPIter) Next(ip net.IP) bool {
	if i.flags&ipIterFlagV6 > 0 {
		if i.flags&ipIterFlagNegative > 0 {
			if i.v6.val.Cmp(i.v6.limit) == -1 {
				return false
			}
			from128(i.v6.val, ip)
			i.v6.val = i.v6.val.Minus(i.v6.incr)
			return true
		}
		if i.v6.val.Cmp(i.v6.limit) != -1 {
			return false
		}
		from128(i.v6.val, ip)
		i.v6.val = i.v6.val.Add(i.v6.incr)
		return true
	}
	if i.flags&ipIterFlagNegative > 0 {
		if i.v4.val < i.v4.limit {
			return false
		}
		from32(i.v4.val, ip)
		i.v4.val -= i.v4.incr
		return true
	}
	if i.v4.val >= i.v4.limit {
		return false
	}
	from32(i.v4.val, ip)
	i.v4.val += i.v4.incr
	return true
}

// IterIP returns an Iter for the given increment over the IPs
func IterIP(ip net.IP, incr int) *IPIter {
	iter := IPIter{}
	if ip.To4() == nil {
		iter.flags |= ipIterFlagV6
	}
	if incr < 0 {
		iter.flags |= ipIterFlagNegative
	}

	switch iter.flags {
	case ipIterFlagV6 | ipIterFlagNegative:
		iter.v6.val, iter.v6.incr, iter.v6.limit = to128(ip), uint128{0, uint64(incr * -1)}, uint128{}
	case ipIterFlagV6:
		iter.v6.val, iter.v6.incr, iter.v6.limit = to128(ip), uint128{0, uint64(incr)}, uint128{maxUint64, maxUint64}
	case ipIterFlagNegative:
		iter.v4.val, iter.v4.incr, iter.v4.limit = to32(ip), uint32(incr*-1), 0
	default:
		iter.v4.val, iter.v4.incr, iter.v4.limit = to32(ip), uint32(incr), maxUint32
	}
	return &iter
}

// NetIter permits iteration over a series of IP networks. It is always start inclusive.
type NetIter struct {
	ips  IPIter
	mask net.IPMask
}

// Next returns true when the underlying pointer has been successfully updated with the next value.
func (n *NetIter) Next(ipN *net.IPNet) bool {
	if !n.ips.Next(ipN.IP) {
		return false
	}
	copy(ipN.Mask, n.mask)
	return true
}

// IterNet returns an iterator for the given increment starting with the provided network
func IterNet(ipNet *net.IPNet, incr int) *NetIter {
	mask := make(net.IPMask, len(ipNet.Mask))
	copy(mask, ipNet.Mask)

	ones, bits := mask.Size()
	suffix := uint(bits - ones)

	iter := IPIter{}
	if ipNet.IP.To4() == nil {
		iter.flags |= ipIterFlagV6
	}
	if incr < 0 {
		iter.flags |= ipIterFlagNegative
	}

	switch iter.flags {
	case ipIterFlagV6 | ipIterFlagNegative:
		iter.v6.val, iter.v6.incr, iter.v6.limit = to128(ipNet.IP), uint128{0, uint64(incr * -1)}.Lsh(suffix), uint128{}
	case ipIterFlagV6:
		iter.v6.val, iter.v6.incr, iter.v6.limit = to128(ipNet.IP), uint128{0, uint64(incr)}.Lsh(suffix), uint128{maxUint64, maxUint64}
	case ipIterFlagNegative:
		iter.v4.val, iter.v4.incr, iter.v4.limit = to32(ipNet.IP), uint32(incr*-1)<<suffix, 0
	default:
		iter.v4.val, iter.v4.incr, iter.v4.limit = to32(ipNet.IP), uint32(incr)<<suffix, maxUint32
	}

	return &NetIter{iter, mask}
}
