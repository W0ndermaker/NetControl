package scanner

import (
	"bufio"
	"net"
	"net/netip"
	"os"
	"strings"
	"time"
)

func poke(ip netip.Addr) {
	conn, err := net.Dial("udp", net.JoinHostPort(ip.String(), "9"))
	if err != nil {
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte{0})
}

func ARPScan(firstIP, lastIP uint32) {
	if lastIP-firstIP >= 2 {
		firstIP, lastIP = firstIP+1, lastIP-1
	}

	for n := firstIP; ; n++ {
		ip := uint32ToAddr(n)

		poke(ip)

		if n == lastIP {
			break
		}
	}

	time.Sleep(time.Second) // дать ядру завершить ARP-запросы
}

func ReadARPTable() (map[string]string, error) {
	f, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	res := map[string]string{}
	sc := bufio.NewScanner(f)
	sc.Scan() // пропускаем заголовок
	for sc.Scan() {
		// IP address | HW type | Flags | HW address | Mask | Device
		fields := strings.Fields(sc.Text())
		if len(fields) < 4 {
			continue
		}
		ip, flags, mac := fields[0], fields[2], fields[3]
		if flags == "0x0" || mac == "00:00:00:00:00:00" {
			continue // неполная запись: устройство не ответило
		}
		res[ip] = mac
	}
	return res, sc.Err()

}

func uint32ToAddr(n uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	})
}
