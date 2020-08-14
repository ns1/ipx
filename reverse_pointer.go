package ipx

import (
	"bytes"
	"fmt"
	"net"
)

// ReversePointer returns the name of the reverse DNS PTR record for the IP address
func ReversePointer(ip net.IP) string {
	var buffer bytes.Buffer

	four := ip.To4()
	if four != nil {
		for i := len(four) - 1; i >= 0; i-- {
			buffer.WriteString(fmt.Sprintf("%d.", four[i]))
		}
		buffer.WriteString("in-addr.arpa")
		return buffer.String()
	}

	for i := len(ip) - 1; i >= 0; i-- {
		b := ip[i]
		buffer.WriteString(fmt.Sprintf("%x.", b&0xF))
		buffer.WriteString(fmt.Sprintf("%x.", b>>4))
	}
	buffer.WriteString("ip6.arpa")
	return buffer.String()
}
