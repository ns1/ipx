package ipx

import (
	"math"
	"math/big"
	"net"
)

var maxUint128 = initMaxUint128()

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
	Next(ipNet net.IPNet) bool
}

// IterIP returns an Iter for the given increment over the IPs
func IterIP(ip net.IP, incr int) IPIter {
	return &deferredIPIter{IPIter: iterIP(ip, incr)}
}

func iterIP(ip net.IP, incr int) IPIter {
	is4 := ip.To4() != nil
	sub := incr < 0
	switch {
	case is4 && sub:
		return newDecrIP4(ip, incr)
	case is4:
		return newIncrIP4(ip, incr)
	case sub:
		return newDecrIP6(ip, big.NewInt(int64(incr)))
	default:
		return newIncrIP6(ip, big.NewInt(int64(incr)))
	}
}

// IterNet returns an iterator for the given increment starting with the provided network
func IterNet(ipNet net.IPNet, incr int) NetIter {
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
			return newDecrIP4(ip, incr<<suffix)
		case is4:
			return newIncrIP4(ip, incr<<suffix)
		case sub:
			incrB := big.NewInt(int64(incr))
			incrB.Lsh(incrB, uint(suffix))
			return newDecrIP6(ip, incrB)
		default:
			incrB := big.NewInt(int64(incr))
			incrB.Lsh(incrB, uint(suffix))
			return newIncrIP6(ip, incrB)
		}
	}()
	return &netIter{IPIter: &deferredIPIter{IPIter: ipIter}, mask: mask}
}

type netIter struct {
	IPIter
	mask net.IPMask
}

func (n *netIter) Next(ipNet net.IPNet) bool {
	if !n.IPIter.Next(ipNet.IP) {
		return false
	}
	copy(ipNet.Mask, n.mask)
	return true
}

type decrIP4 struct {
	v, decr uint32
}

func (d *decrIP4) Next(ip net.IP) bool {
	if d.v < d.decr {
		return false
	}
	d.v -= d.decr
	from32(d.v, ip)
	return true
}

type incrIP4 struct {
	v, incr uint32
}

func (i *incrIP4) Next(ip net.IP) bool {
	if math.MaxUint32-i.incr < i.v {
		return false
	}
	i.v += i.incr
	from32(i.v, ip)
	return true
}

type incrIP6 struct {
	limit, v, incr *big.Int
}

func (i *incrIP6) Next(ip net.IP) bool {
	if i.limit.Cmp(i.v) != 1 {
		return false
	}
	i.v.Add(i.v, i.incr)
	from128(i.v, ip)
	return true
}

type decrIP6 struct {
	v, decr *big.Int
}

func (d *decrIP6) Next(ip net.IP) bool {
	if d.decr.Cmp(d.v) == 1 {
		return false
	}
	d.v.Sub(d.v, d.decr)
	from128(d.v, ip)
	return true
}

func newDecrIP4(ip net.IP, incr int) *decrIP4 {
	return &decrIP4{
		v:    to32(ip),
		decr: uint32(incr * -1),
	}
}

func newIncrIP4(ip net.IP, incr int) *incrIP4 {
	return &incrIP4{
		v:    to32(ip),
		incr: uint32(incr),
	}
}

func newDecrIP6(ip net.IP, incr *big.Int) *decrIP6 {
	decr := incr.Neg(incr)
	return &decrIP6{
		v:    to128(ip),
		decr: decr,
	}
}

func newIncrIP6(ip net.IP, incr *big.Int) *incrIP6 {
	limit := new(big.Int).Sub(maxUint128, incr)
	return &incrIP6{
		v:     to128(ip),
		incr:  incr,
		limit: limit,
	}
}

type deferredIPIter struct {
	first bool
	IPIter
}

func (d *deferredIPIter) Next(ip net.IP) bool {
	if !d.first {
		d.first = true
		return true
	}
	return d.IPIter.Next(ip)
}

func initMaxUint128() *big.Int {
	i := new(big.Int).SetUint64(math.MaxUint64)
	i.Lsh(i, 64)
	i.Or(i, new(big.Int).SetUint64(math.MaxUint64))
	return i
}
