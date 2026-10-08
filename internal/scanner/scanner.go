package scanner

import (
	"fmt"
	"net"
)

type NetworkAddress struct {
	IP   net.IP
	Mask int
}

func getIPv4Addrs(iface net.Interface) ([]NetworkAddress, error) {

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении адресов из сетевого интерфейса: %w", err)
	}

	netAddrs := make([]NetworkAddress, 0, len(addrs))

	for _, addr := range addrs {
		network, ok := addr.(*net.IPNet)

		if !ok {
			continue
		}

		if network.IP.To4() == nil {
			continue
		}

		ones, _ := network.Mask.Size()

		netAddrs = append(netAddrs, NetworkAddress{
			network.IP,
			ones,
		})
	}

	return netAddrs, nil
}

func NetScan() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("Ошибка при получении сетевых интерфейсов: %w", err)
	}

	for _, iface := range interfaces {

		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		ipv4Addrs, err := getIPv4Addrs(iface)
		if err != nil {
			return err
		}

		if len(ipv4Addrs) == 0 {
			continue
		}

		//fmt.Printf("%v.\n\t", iface.Index)
		//fmt.Printf("Максимальный размер сетевого пакета: %v\n\t", iface.MTU)
		fmt.Printf("\n\tИмя сетевого интерфейса: %v\n\t", iface.Name)
		fmt.Printf("MAC-адрес: %v\n\t", iface.HardwareAddr)
		fmt.Printf("Флаги сетевого интерфейса: %v\n", iface.Flags)

		fmt.Printf("\tIPv4: \n")
		for _, ipv4Addr := range ipv4Addrs {
			fmt.Printf("\t\t%v/%v\n", ipv4Addr.IP, ipv4Addr.Mask)
		}
		fmt.Println()

	}

	return nil
}
