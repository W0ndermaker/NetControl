package scanner

import (
	"fmt"
	"net"
)

type NetworkAddress struct {
	IP        net.IP
	Mask      net.IPMask
	NetworkIP net.IP
}

// Сканируем сетевой интерфейс для получения данных сети
func getInterfaceInfo(iface net.Interface) ([]NetworkAddress, error) {

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

		mask := network.Mask
		_, bits := mask.Size()

		if bits != 32 {
			continue
		}

		netAddrs = append(netAddrs, NetworkAddress{
			IP:   network.IP,
			Mask: network.Mask,
			// Накидывается маска сети на локальный IP-адрес для вычисления IP-адреса сети
			NetworkIP: network.IP.Mask(mask),
		})
	}

	return netAddrs, nil
}

// Сканируем все доступные сетевые интерфейсы для текущего устройства
func ScanNetworkInterfaces() error {
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

		ipv4Addrs, err := getInterfaceInfo(iface)
		if err != nil {
			return err
		}

		if len(ipv4Addrs) == 0 {
			continue
		}

		showInterfaces(iface, ipv4Addrs)
	}

	return nil
}
