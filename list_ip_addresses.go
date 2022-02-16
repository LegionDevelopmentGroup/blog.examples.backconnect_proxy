package main

import (
	"fmt"
	"net"
)

func main() {
	targetInterfaceMAC := "9c:b6:d0:dd:5e:19"
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var activeIPList []string
	for _, netInterface := range netInterfaces {
		if netInterface.HardwareAddr.String() != targetInterfaceMAC {
			continue
		}
		addresses, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		for _, address := range addresses {
			if val, ok := address.(*net.IPNet); ok {
				if val.IP.To4() != nil { // Filters out IPv6 addresses
					activeIPList = append(activeIPList, val.IP.String())
				}
			}
		}
	}

	fmt.Printf("MAC %s has the following addresses:\n", targetInterfaceMAC)
	for i, ip := range activeIPList {
		fmt.Printf("%d - %s\n", i, ip)
	}
}
