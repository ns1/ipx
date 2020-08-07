package ipx

// largely cribbed from https://github.com/davidminor/uint128
type uint128 struct {
	H, L uint64
}

func (u *uint128) And(other uint128) {
	u.H &= other.H
	u.L &= other.L
}

func (u *uint128) Or(other uint128) {
	u.H |= other.H
	u.L |= other.L
}

func (u uint128) Cmp(other uint128) int {
	switch {
	case u.H > other.H:
		return 1
	case u.H < other.H,
		u.L < other.L:
		return -1
	case u.L > other.L:
		return 1
	default:
		return 0
	}
}

func (u *uint128) Add(addend uint128) {
	old := u.L
	u.H += addend.H
	u.L += addend.L
	if u.L < old { // wrapped
		u.H += 1
	}
}

func (u *uint128) Minus(addend uint128) {
	old := u.L
	u.H -= addend.H
	u.L -= addend.L
	if u.L > old { // wrapped
		u.H -= 1
	}
}

func (u *uint128) Lsh(bits uint) {
	switch {
	case bits >= 128:
		u.H, u.L = 0, 0
	case bits >= 64:
		u.H, u.L = u.L<<(bits-64), 0
	default:
		u.H <<= bits
		u.H |= u.L >> (64 - bits) // set top with bits that cross from bottom
		u.L <<= bits
	}
}

func (u *uint128) Rsh(bits uint) {
	switch {
	case bits >= 128:
		u.H, u.L = 0, 0
	case bits >= 64:
		u.H, u.L = 0, u.H>>(bits-64)
	default:
		u.L >>= bits
		u.L |= u.H << (64 - bits) // set bottom with bits that cross from top
		u.H >>= bits
	}
}

func to128(ip []byte) uint128 {
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

func from128(u uint128, ip []byte) {
	ip[0] = byte(u.H >> 56)
	ip[1] = byte(u.H >> 48)
	ip[2] = byte(u.H >> 40)
	ip[3] = byte(u.H >> 32)
	ip[4] = byte(u.H >> 24)
	ip[5] = byte(u.H >> 16)
	ip[6] = byte(u.H >> 8)
	ip[7] = byte(u.H)
	ip[8] = byte(u.L >> 56)
	ip[9] = byte(u.L >> 48)
	ip[10] = byte(u.L >> 40)
	ip[11] = byte(u.L >> 32)
	ip[12] = byte(u.L >> 24)
	ip[13] = byte(u.L >> 16)
	ip[14] = byte(u.L >> 8)
	ip[15] = byte(u.L)
}
