package scanner

import "net"

func ipToUint32(ip net.IP) uint32 {
	return uint32(ip[0])<<24 |
		uint32(ip[1])<<16 |
		uint32(ip[2])<<8 |
		uint32(ip[3])
}

func uint32ToIp(ipNum uint32) net.IP {
	return net.IP{
		byte(ipNum >> 24),
		byte(ipNum >> 16),
		byte(ipNum >> 8),
		byte(ipNum),
	}
}
