package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const IP_PORTION = 5

func showInterfaces(iface net.Interface, ipv4Addrs []NetworkAddress) {
	//fmt.Printf("%v.\n\t", iface.Index)
	//fmt.Printf("Максимальный размер сетевого пакета: %v\n\t", iface.MTU)
	fmt.Printf("\n\tИмя сетевого интерфейса: %v\n\t", iface.Name)
	fmt.Printf("MAC-адрес: %v\n\t", iface.HardwareAddr)
	fmt.Printf("Флаги сетевого интерфейса: %v\n", iface.Flags)

	fmt.Printf("\tIPv4: \n")
	for i, ipv4Addr := range ipv4Addrs {
		if i != 0 {
			fmt.Println()
		}

		ipNet := net.IPNet{
			IP:   ipv4Addr.NetworkIP,
			Mask: ipv4Addr.Mask,
		}

		fmt.Printf("\t\tIP: %v\n", ipv4Addr.IP)
		fmt.Printf("\t\tМаска: %v\n", net.IP(ipv4Addr.Mask).To4())
		fmt.Printf("\t\tАдрес сети: %v\n", ipNet.String())
	}

}

func showIPRange(networkIP net.IP, broadcast net.IP) error {
	fmt.Printf("Адрес сети: %v\n", networkIP)
	fmt.Printf("Broadcast: %v\n\n", broadcast)

	networkNumber := ipToUint32(networkIP)
	broadcastNumber := ipToUint32(broadcast)

	if broadcastNumber-networkNumber < 2 {
		return nil
	}
	
	firstHostIP := uint32ToIp(networkNumber + 1)
	lastHostIP := uint32ToIp(broadcastNumber - 1)

	fmt.Println("Первый адрес хоста:", firstHostIP)
	fmt.Println("Последний адрес хоста:", lastHostIP)

	count := 0
	fmt.Printf("\nIP-адреса:\n")

	reader := bufio.NewReader(os.Stdin)

	for i := networkNumber + 1; i <= broadcastNumber-1; i++ {
		fmt.Println("\t", uint32ToIp(i))
		count++

		if count == IP_PORTION && i < broadcastNumber-1 {
			con, err := getCommandContinue(reader)
			if err != nil {
				return err
			}

			if !con {
				break
			}

			count = 0
		}
	}
	return nil
}

func getCommandContinue(r *bufio.Reader) (bool, error) {

	fmt.Println("Нажмите Enter для продолжения или q для выхода")

	input, err := r.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Выход")
		return false, nil
	}

	return true, nil
}
