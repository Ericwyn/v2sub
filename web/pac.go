package web

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func renderPacJs(ipAddress string, hPort string, sPort string) string {
	if ipAddress == "" {
		ipAddress = getIp(false)
	}
	if sPort == "" {
		sPort = "1080"
	}
	if hPort == "" {
		hPort = "1081"
	}
	res := strings.Replace(pacJs, "SOCKS5 127.0.0.1:1080", "PROXY ${ipAddress}:${hPort}; SOCKS5 ${ipAddress}:${sPort}; DIRECT", 1)
	res = strings.Replace(res, "${ipAddress}", ipAddress, -1)
	res = strings.Replace(res, "${hPort}", hPort, -1)
	res = strings.Replace(res, "${sPort}", sPort, -1)

	return res
}

func getIp(isIpv6 bool) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resIp := ""
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ipParse := ipnet.IP.String()

			isIp := false
			if isIpv6 {
				isIp = IsIPv6(ipParse)
			} else {
				isIp = IsIPv4(ipParse)
			}

			if ipParse != "" && isIp {
				if resIp == "" {
					resIp = ipParse
				} else {
					// 如果当前 ip 已经存在，但是新的 ip 并不是本地地址的话，可以覆盖掉
					if !isLocalAddress(ipParse) {
						resIp = ipParse
					}
				}
			}
		}
	}
	return resIp
}

func getIpv6() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			//if ipnet.IP.To4() != nil {
			//	return ipnet.IP.String()
			//}
			ipParse := ipnet.IP.String()
			if ipParse != "" && IsIPv6(ipParse) {
				return ipParse
			}
		}
	}
	return ""
}

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

// 判断是否为本地地址
func isLocalAddress(address string) bool {
	return strings.Index(address, "fe80") == 0 ||
		strings.Index(address, "127") == 0 ||
		strings.Index(address, "192") == 0
}
