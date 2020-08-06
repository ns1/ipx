package ipx

import "net"

// largely cribbed from https://github.com/davidminor/uint128
type uint128 [2]uint64

func (u *uint128) And(other uint128) {
	u[0] &= other[0]
	u[1] &= other[1]
}

func (u *uint128) Or(other uint128) {
	u[0] |= other[0]
	u[1] |= other[1]
}

func (u uint128) Cmp(other uint128) int {
	switch {
	case u[0] > other[0]:
		return 1
	case u[0] < other[0],
		u[1] < other[1]:
		return -1
	case u[1] > other[1]:
		return 1
	default:
		return 0
	}
}

func (u *uint128) Add(addend uint128) {
	old := u[1]
	u[0] += addend[0]
	u[1] += addend[1]
	if u[1] < old { // wrapped
		u[0] += 1
	}
}

func (u *uint128) Minus(addend uint128) {
	old := u[1]
	u[0] -= addend[0]
	u[1] -= addend[1]
	if u[1] > old { // wrapped
		u[0] -= 1
	}
}

func (u *uint128) Lsh(bits uint) {
	switch {
	case bits >= 128:
		u[0], u[1] = 0, 0
	case bits >= 64:
		u[0], u[1] = u[1]<<(bits-64), 0
	default:
		u[0] <<= bits
		u[0] |= u[1] >> (64 - bits) // set top with bits that cross from bottom
		u[1] <<= bits
	}
}

func (u *uint128) Rsh(bits uint) {
	switch {
	case bits >= 128:
		u[0], u[1] = 0, 0
	case bits >= 64:
		u[0], u[1] = 0, u[0]>>(bits-64)
	default:
		u[1] >>= bits
		u[1] |= u[0] << (64 - bits) // set bottom with bits that cross from top
		u[0] >>= bits
	}
}

func toUint128(ip net.IP) uint128 {
	return uint128{
		uint64(ip[0])<<56 |
			uint64(ip[1])<<48 |
			uint64(ip[2])<<40 |
			uint64(ip[3])<<32 |
			uint64(ip[4])<<24 |
			uint64(ip[5])<<16 |
			uint64(ip[6])<<8 |
			uint64(ip[7]),
		uint64(ip[8])<<56 |
			uint64(ip[9])<<48 |
			uint64(ip[10])<<40 |
			uint64(ip[11])<<32 |
			uint64(ip[12])<<24 |
			uint64(ip[13])<<16 |
			uint64(ip[14])<<8 |
			uint64(ip[15]),
	}
}

func fromUint128(u uint128, ip net.IP) {
	ip[0] = byte(u[0] >> 56)
	ip[1] = byte(u[0] >> 48)
	ip[2] = byte(u[0] >> 40)
	ip[3] = byte(u[0] >> 32)
	ip[4] = byte(u[0] >> 24)
	ip[5] = byte(u[0] >> 16)
	ip[6] = byte(u[0] >> 8)
	ip[7] = byte(u[0])
	ip[8] = byte(u[1] >> 56)
	ip[9] = byte(u[1] >> 48)
	ip[10] = byte(u[1] >> 40)
	ip[11] = byte(u[1] >> 32)
	ip[12] = byte(u[1] >> 24)
	ip[13] = byte(u[1] >> 16)
	ip[14] = byte(u[1] >> 8)
	ip[15] = byte(u[1])
}
