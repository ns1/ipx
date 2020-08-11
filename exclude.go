package ipx

import "net"

// Exclude returns a list of networks representing the address block when `b` is removed from `a`.
func Exclude(a, b *net.IPNet) []*net.IPNet {
	four := a.IP.To4() != nil
	if four != (b.IP.To4() != nil) || !IsSubnet(a, b) {
		return []*net.IPNet{a}
	}
	if a.IP.To4() != nil {
		return exclude4(newIP4Net(a), newIP4Net(b))
	}
	return exclude6(newIP6Net(a), newIP6Net(b))
}

func exclude4(a, b ip4Net) []*net.IPNet {
	subs := make([]*net.IPNet, 0, a.prefix-b.prefix)

	s1, s2 := a.subnets()
	for s1 != b && s2 != b {
		if b.subnetOf(s1) {
			subs = append(subs, s2.asNet())
			s1, s2 = s1.subnets()
			continue
		}
		subs = append(subs, s1.asNet())
		s1, s2 = s2.subnets()
	}
	if s1 == b {
		subs = append(subs, s2.asNet())
	} else {
		subs = append(subs, s1.asNet())
	}
	return subs
}

func exclude6(a, b ip6Net) []*net.IPNet {
	subs := make([]*net.IPNet, 0, a.prefix-b.prefix)

	s1, s2 := a.subnets()
	for s1 != b && s2 != b {
		if b.subnetOf(s1) {
			subs = append(subs, s2.asNet())
			s1, s2 = s1.subnets()
			continue
		}
		subs = append(subs, s1.asNet())
		s1, s2 = s2.subnets()
	}
	if s1 == b {
		subs = append(subs, s2.asNet())
	} else {
		subs = append(subs, s1.asNet())
	}
	return subs
}
