package ipx

import (
	"fmt"
	"net"
)

// Split splits a subnet into smaller subnets according to the new prefix provided.
func Split(ipNet *net.IPNet, newPrefix int) NetIter {
	ones, bits := ipNet.Mask.Size()
	if ones > newPrefix || newPrefix > bits {
		panic(fmt.Errorf("must be in [%v, %v] but got %v", ones, bits, newPrefix))
	}
	if ipNet.IP.To4() != nil {
		ip := to32(ipNet.IP)
		return &incrIP4Net{
			incrIP4{
				ip,
				1 << (bits - newPrefix),
				ip | (1<<(bits-ones) - 1),
			},
			((1 << newPrefix) - 1) << (bits - newPrefix),
		}
	}

	ip := to128(ipNet.IP)

	incr := uint128{0, 1}
	incr.Lsh(uint(bits - newPrefix))

	broadCast := uint128{0, 1}
	broadCast.Lsh(uint(bits - ones))
	broadCast.Minus(uint128{0, 1})
	broadCast.Or(ip)

	mask := uint128{0, 1}
	mask.Lsh(uint(newPrefix))
	mask.Minus(uint128{0, 1})
	mask.Lsh(uint(bits - newPrefix))

	return &incrIP6Net{incrIP6{ip, incr, broadCast}, mask}
}

// Addresses returns all of the addresses within a network.
func Addresses(ipNet *net.IPNet) IPIter {
	ones, bits := ipNet.Mask.Size()
	if ipNet.IP.To4() != nil {
		ip := to32(ipNet.IP)
		return &incrIP4{ip, 1, ip + (1 << (bits - ones))}
	}
	ip := to128(ipNet.IP)

	addend := uint128{0, 1}
	addend.Lsh(uint(bits - ones))

	limit := ip
	limit.Add(addend)

	return &incrIP6{ip, uint128{0, 1}, limit}
}

// Hosts returns all of the usable addresses within a network except the network itself address and the broadcast address
func Hosts(ipNet *net.IPNet) IPIter {
	ones, bits := ipNet.Mask.Size()
	if ipNet.IP.To4() != nil {
		ip := to32(ipNet.IP) + 1
		return &incrIP4{ip, 1, ip + (1 << (bits - ones)) - 2}
	}

	ip := to128(ipNet.IP)
	ip.Add(uint128{0, 1})

	addend := uint128{0, 1}
	addend.Lsh(uint(bits - ones))
	addend.Minus(uint128{0, 2})

	limit := ip
	limit.Add(addend)

	return &incrIP6{ip, uint128{0, 1}, limit}
}
