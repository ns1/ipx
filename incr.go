package ipx

import (
	"errors"
	"net"
)

// IncrIP returns the next IP
func IncrIP(ip net.IP, incr int) {
	if ip == nil {
		panic(errors.New("IP cannot be nil"))
	}
	if ip.To4() != nil {
		n := to32(ip)
		if incr >= 0 {
			n += uint32(incr)
		} else {
			n -= uint32(incr)
		}
		from32(n, ip)
		return
	}

	// ipv6
	u := toUint128(ip)
	if incr >= 0 {
		u.Add(uint128{0, uint64(incr)})
	} else {
		u.Minus(uint128{0, uint64(incr * -1)})
	}
	fromUint128(u, ip)
}

// IncrNet steps to the next net of the same mask
func IncrNet(ipNet *net.IPNet, incr int) {
	if ipNet.IP == nil {
		panic(errors.New("IP cannot be nil"))
	}
	if ipNet.Mask == nil {
		panic(errors.New("mask cannot be nil"))
	}
	if ipNet.IP.To4() != nil {
		n := to32(ipNet.IP)
		ones, bits := ipNet.Mask.Size()
		suffix := uint32(bits - ones)
		n >>= suffix
		if incr >= 0 {
			n += uint32(incr)
		} else {
			n -= uint32(incr)
		}
		from32(n<<suffix, ipNet.IP)
		return
	}

	b := toUint128(ipNet.IP)

	ones, bits := ipNet.Mask.Size()
	suffix := uint(bits - ones)

	b.Rsh(suffix)
	if incr >= 0 {
		b.Add(uint128{0, uint64(incr)})
	} else {
		b.Minus(uint128{0, uint64(incr * -1)})
	}
	b.Lsh(suffix)

	fromUint128(b, ipNet.IP)
}

func to32(ip net.IP) uint32 {
	l := len(ip)
	return uint32(ip[l-4])<<24 |
		uint32(ip[l-3])<<16 |
		uint32(ip[l-2])<<8 |
		uint32(ip[l-1])
}

func from32(n uint32, ip net.IP) {
	l := len(ip)
	ip[l-4] = uint8(n >> 24)
	ip[l-3] = uint8(n >> 16)
	ip[l-2] = uint8(n >> 8)
	ip[l-1] = uint8(n)
}
