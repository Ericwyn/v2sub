package web

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// 参考 pac 文档
// 因为 v2sub 的路由文件其实已经很完善了
// android 手机端可以直接设置代理为 http 端口
// 但是如果 v2sub 的代理关闭了的话, 那就连普通网页都上不了(不会自动切换)
// pac 脚本的话就可以解决这个问题
// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
const pacJs = `
function FindProxyForURL(url, host) {
    // If the protocol or URL matches, send direct.
    if (url.substring(0, 4)=="ftp:")
        return "DIRECT";
     
    // If the requested website is hosted within the internal network, send direct.
    if (isPlainHostName(host) ||
        shExpMatch(host, "*.local") ||
        isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
        isInNet(dnsResolve(host), "172.16.0.0",  "255.240.0.0") ||
        isInNet(dnsResolve(host), "192.168.0.0",  "255.255.0.0") ||
        isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0")) {
        return "DIRECT";
    }
          
    // DEFAULT RULE: All other traffic, use below proxies, in fail-over order.
    return "PROXY ${ipAddress}:${hPort}; SOCKS5 ${ipAddress}:${sPort}; DIRECT";
}
`

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
	res := strings.Replace(pacJs, "${ipAddress}", ipAddress, -1)
	res = strings.Replace(res, "${hPort}", hPort, -1)
	res = strings.Replace(res, "${sPort}", sPort, -1)

	fmt.Println(res)

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
