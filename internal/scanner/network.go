package scanner

import (
	"errors"
	"fmt"
	"math"
	"net"
)

func ScanNetwork(RawNetworkIP string) error {

	_, ipNet, err := net.ParseCIDR(RawNetworkIP)
	if err != nil {
		return fmt.Errorf("Ошибка при вводе IP: %v", err)
	}

	return generateIPRange(ipNet)
}

func generateIPRange(ipNet *net.IPNet) error {
	ipCount := networkSize(ipNet)
	fmt.Println("Количество IP адресов в сети: ", ipCount)

	networkIP := ipNet.IP.To4()
	if networkIP == nil {
		return errors.New("Нужен IPv4")
	}
	mask := ipNet.Mask
	fmt.Println("Маска: ", net.IP(mask))

	broadcast := make(net.IP, 4)

	for i := range len(networkIP) {
		broadcast[i] = networkIP[i] | ^mask[i]
	}

	err := showIPRange(networkIP, broadcast)
	if err != nil {
		return nil
	}
	return nil
}

func networkSize(ipNet *net.IPNet) int {
	ones, bits := ipNet.Mask.Size()
	return int(math.Pow(2, float64(bits-ones)))
}
