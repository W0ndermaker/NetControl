package scanner

import (
	"errors"
	"fmt"
	"net"
)

const maxAddresses uint64 = 65536

func ScanNetwork(RawNetworkIP string) error {

	_, ipNet, err := net.ParseCIDR(RawNetworkIP)
	if err != nil {
		return fmt.Errorf("Ошибка при вводе IP: %v", err)
	}

	return generateIPRange(ipNet)
}

func generateIPRange(ipNet *net.IPNet) error {
	ipCount, err := networkSize(ipNet)
	if err != nil {
		return err
	}
	fmt.Println("Количество IP адресов в сети: ", ipCount)

	networkIP := ipNet.IP.To4()
	if networkIP == nil {
		return errors.New("Нужен IPv4")
	}
	mask := ipNet.Mask
	fmt.Println("Маска: ", mask)

	broadcast := make(net.IP, 4)

	for i := range len(networkIP) {
		broadcast[i] = networkIP[i] | ^mask[i]
	}

	err = showIPRange(networkIP, broadcast)
	if err != nil {
		return err
	}
	return nil
}

func networkSize(ipNet *net.IPNet) (uint64, error) {
	ones, bits := ipNet.Mask.Size()

	if bits != 32 {
		return 0, fmt.Errorf("Поддерживается только IPv4!")
	}

	totalAddrs := uint64(1) << uint(32-ones)

	if totalAddrs > maxAddresses {
		return 0, fmt.Errorf("Сеть слишком большая %d адресов, максимум %d",
			totalAddrs,
			maxAddresses)
	}

	return totalAddrs, nil
}
