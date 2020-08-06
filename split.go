package ipx

import (
	"fmt"
	"net"
)

// Split
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
