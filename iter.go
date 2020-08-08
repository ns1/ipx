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

	ip net.IP
}

// IP returns the most recent IP; the underlying type may be modified on later calls to `Next`.
// It does no allocation.
func (i *IPIter) IP() net.IP {
	return i.ip
}

// Next returns true when the underlying pointer has been successfully updated with the next value.
func (i *IPIter) Next() bool {
	if i.flags&ipIterFlagV6 > 0 {
		if i.flags&ipIterFlagNegative > 0 {
			if i.v6.val.Cmp(i.v6.limit) == -1 {
				return false
			}
			from128(i.v6.val, i.ip)
			i.v6.val = i.v6.val.Minus(i.v6.incr)
			return true
		}
		if i.v6.val.Cmp(i.v6.limit) != -1 {
			return false
		}
		from128(i.v6.val, i.ip)
		i.v6.val = i.v6.val.Add(i.v6.incr)
		return true
	}
	if i.flags&ipIterFlagNegative > 0 {
		if i.v4.val < i.v4.limit {
			return false
		}
		from32(i.v4.val, i.ip)
		i.v4.val -= i.v4.incr
		return true
	}
	if i.v4.val >= i.v4.limit {
		return false
	}
	from32(i.v4.val, i.ip)
	i.v4.val += i.v4.incr
	return true
}

// IterIP returns an Iter for the given increment over the IPs
func IterIP(ip net.IP, incr int) *IPIter {
	if ip.To4() != nil {
		var (
			uIncr uint32
			limit uint32
		)
		if incr < 0 {
			uIncr = uint32(incr * -1)
			limit = 0
		} else {
			uIncr = uint32(incr)
			limit = maxUint32
		}
		return iterIPv4(to32(ip), uIncr, limit)
	}
	var (
		uIncr uint128
		limit uint128
	)
	if incr < 0 {
		uIncr = uint128{0, uint64(incr * -1)}
		limit = uint128{}
	} else {
		uIncr = uint128{0, uint64(incr)}
		limit = uint128{maxUint64, maxUint64}
	}
	return iterIPv6(to128(ip), uIncr, limit, )
}

func iterIPv4(val, incr, limit uint32) *IPIter {
	iter := IPIter{ip: make(net.IP, len(net.IPv4zero)), v4: v4IPIter{val, incr, limit}}
	copy(iter.ip, net.IPv4zero)
	if limit < val {
		iter.flags |= ipIterFlagNegative
	}
	return &iter
}

func iterIPv6(val, incr, limit uint128) *IPIter {
	iter := IPIter{
		ip:    make(net.IP, len(net.IPv6zero)),
		v6:    v6IPIter{val, incr, limit},
		flags: ipIterFlagV6,
	}
	copy(iter.ip, net.IPv6zero)
	if limit.Cmp(val) == -1 {
		iter.flags |= ipIterFlagNegative
	}
	return &iter
}

// NetIter permits iteration over a series of IP networks. It is always start inclusive.
type NetIter struct {
	ips IPIter
	net *net.IPNet
}

// Net returns the most recent IPNet; the underlying type may be modified on later calls to `Next`.
// It does no allocation.
func (n *NetIter) Net() *net.IPNet {
	n.net.IP = n.ips.IP()
	return n.net
}

// Next returns true when the underlying pointer has been successfully updated with the next value.
func (n *NetIter) Next() bool {
	return n.ips.Next()
}

// IterNet returns an iterator for the given increment starting with the provided network
func IterNet(ipNet *net.IPNet, incr int) *NetIter {
	mask := make(net.IPMask, len(ipNet.Mask))
	copy(mask, ipNet.Mask)

	ones, bits := mask.Size()
	suffix := uint(bits - ones)

	if ipNet.IP.To4() != nil {
		if incr < 0 {
			return &NetIter{
				*iterIPv4(
					to32(ipNet.IP),
					uint32(incr*-1)<<suffix,
					0,
				),
				&net.IPNet{Mask: mask},
			}
		}
		return &NetIter{
			*iterIPv4(
				to32(ipNet.IP),
				uint32(incr)<<suffix,
				maxUint32,
			),
			&net.IPNet{Mask: mask},
		}
	}
	if incr < 0 {
		return &NetIter{
			*iterIPv6(to128(ipNet.IP), uint128{0, uint64(incr * -1)}.Lsh(suffix), uint128{}, ),
			&net.IPNet{Mask: mask},
		}
	}
	return &NetIter{
		*iterIPv6(to128(ipNet.IP), uint128{0, uint64(incr)}.Lsh(suffix), uint128{maxUint64, maxUint64}, ),
		&net.IPNet{Mask: mask},
	}
}
