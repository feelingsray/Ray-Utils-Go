package nethelper

import (
	"net"
)

func GetIPAddressByName(name string) string {
	addr, _ := net.InterfaceByName(name)
	addrs, _ := addr.Addrs()
	data := ""
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				data = ipnet.IP.String()
			}

		}
	}
	return data
}
