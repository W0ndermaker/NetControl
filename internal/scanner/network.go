package scanner

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func ipv4Network(ip net.IP, mask net.IPMask) (start, end uint32, err error) {
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, 0, errors.New("нужен IPv4-адрес")
	}

	ones, bits := mask.Size()
	if bits == 0 {
		return 0, 0, errors.New("маска не непрерывная или неверной длины")
	}

	if bits == 128 {
		ones -= 96
	}
	if ones < 0 || ones > 32 {
		return 0, 0, errors.New("некорректная маска")
	}

	addr := binary.BigEndian.Uint32(ip4)
	netmask := ^uint32(0) << (32 - ones)

	start = addr & netmask
	// broadcast
	end = addr | ^netmask
	return start, end, nil
}

func selectAddress(addrs []NetworkAddress) (NetworkAddress, error) {
	if len(addrs) == 1 {
		return addrs[0], nil
	}

	fmt.Println("Выберите, какой IP-адрес сканировать:")
	for i, addr := range addrs {
		fmt.Printf("%d --- %v\n", i+1, addr)
	}

	var choice int
	if _, err := fmt.Scan(&choice); err != nil {
		return NetworkAddress{}, fmt.Errorf("некорректный ввод: %w", err)
	}
	if choice < 1 || choice > len(addrs) {
		return NetworkAddress{}, fmt.Errorf("номер должен быть от 1 до %d", len(addrs))
	}
	return addrs[choice-1], nil
}

func ScanNetwork(interfaceName string) error {

	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return fmt.Errorf("интерфейс %q: %w", interfaceName, err)
	}

	addrs, err := getInterfaceInfo(*iface)
	if err != nil {
		return fmt.Errorf("Ошибка адреса: %w", err)
	}

	if len(addrs) == 0 {
		return fmt.Errorf("у интерфейса %s нет IPv4-адресов", iface.Name)
	}

	address, err := selectAddress(addrs)
	if err != nil {
		return err
	}

	start, end, err := ipv4Network(address.IP, address.Mask)
	if err != nil {
		return err
	}

	fmt.Println("IP адрес: ", address.IP)

	showIPinfo(start, end)

	ARPScan(start, end)
	table, err := ReadARPTable()
	if err != nil {
		fmt.Println("Не удалось прочитать ARP-таблицу:", err)
		return err
	}

	fmt.Println("\nКоличество занятых IP-адресов: ", len(table))
	fmt.Println("IP\t\t\tMAC")
	for ip, mac := range table {
		fmt.Printf("%s\t\t%s\n", ip, mac)
	}

	return nil

}
