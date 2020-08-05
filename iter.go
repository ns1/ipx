package ipx

import (
	"math"
	"math/big"
	"net"
)

var maxUint128 = initMaxUint128()

// IterIP returns an IPIter for the given increment
func IterIP(incr int) IPIter {
	return &deferredIter{incr: incr}
}

// IPIter permits iteration over a series of ips.
type IPIter interface {
	// Next returns true when the underlying pointer has been successfully updated
	// with the next value.
	Next(ip net.IP) bool
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

type deferredIter struct {
	incr int
	IPIter
}

func (d *deferredIter) Next(ip net.IP) bool {
	if d.IPIter != nil {
		return d.IPIter.Next(ip)
	}
	d.IPIter = iterIP(ip, d.incr)
	return true
}

func iterIP(ip net.IP, incr int) IPIter {
	is4 := ip.To4() != nil
	sub := incr < 0
	switch {
	case is4 && sub:
		return &decrIP4{
			v:    to32(ip),
			decr: uint32(incr * -1),
		}
	case is4:
		return &incrIP4{
			v:    to32(ip),
			incr: uint32(incr),
		}
	case sub:
		decr := big.NewInt(int64(incr * -1))
		return &decrIP6{
			v:    to128(ip),
			decr: decr,
		}
	default:
		incr := big.NewInt(int64(incr))
		limit := new(big.Int).Sub(maxUint128, incr)
		return &incrIP6{
			v:     to128(ip),
			incr:  incr,
			limit: limit,
		}
	}
}

func initMaxUint128() *big.Int {
	i := new(big.Int).SetUint64(math.MaxUint64)
	i.Lsh(i, 64)
	i.Or(i, new(big.Int).SetUint64(math.MaxUint64))
	return i
}
