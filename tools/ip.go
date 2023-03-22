package tools

import "net"

func GetIPAddressByName(name string) string {
	addr, _ := net.InterfaceByName(name)
	addrs, _ := addr.Addrs()
	data := ""
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				data = ipNet.IP.String()
			}
		}
	}
	return data
}
