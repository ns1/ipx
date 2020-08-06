package ipx

import (
	"fmt"
	"net"
)

// Split splits a subnet into smaller subnets according to the new prefix provided.
func Split(ipNet net.IPNet, newPrefix int) NetIter {
	ones, bits := ipNet.Mask.Size()
	if ones > newPrefix || newPrefix > bits {
		panic(fmt.Errorf("must be in [%v, %v] but got %v", ones, bits, newPrefix))
	}
	return &limitNetIter{
		iterNet(ipNet.IP, net.CIDRMask(newPrefix, bits), 1),
		1 << uint(newPrefix-ones),
	}
}

type limitNetIter struct {
	NetIter
	rem int
}

func (l *limitNetIter) Next(ipNet net.IPNet) bool {
	if l.rem == 0 {
		return false
	}
	l.rem--
	return l.NetIter.Next(ipNet)
}

// Addresses returns all of the addresses within a network.
func Addresses(ipNet net.IPNet) IPIter {
	ones, bits := ipNet.Mask.Size()
	return &limitIPIter{
		&includeIPIter{
			IPIter: iterIP(ipNet.IP, 1),
			ip:     ipNet.IP,
		},
		1 << (bits - ones),
	}
}

// Hosts returns all of the usable addresses within a network except the network itself address and the broadcast address
func Hosts(ipNet net.IPNet) IPIter {
	ones, bits := ipNet.Mask.Size()
	return &limitIPIter{
		iterIP(ipNet.IP, 1),
		(1 << (bits - ones)) - 2,
	}
}

type limitIPIter struct {
	IPIter
	rem int
}

func (l *limitIPIter) Next(ip net.IP) bool {
	if l.rem == 0 {
		return false
	}
	l.rem--
	return l.IPIter.Next(ip)
}

type includeIPIter struct {
	IPIter
	ip net.IP
}

func (i *includeIPIter) Next(ip net.IP) bool {
	if i.ip != nil {
		copy(ip, i.ip)
		i.ip = nil
		return true
	}
	return i.IPIter.Next(ip)
}
