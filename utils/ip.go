package utils

import (
	"net"
	"regexp"
	"donniezhangzq/goraft/constant"
)

func GetLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && CheckIpAddress(ipnet.IP.String()) {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", constant.ErrGetLocalIpFailed
}

func CheckIpAddress(ip string) bool {
	result, _ := regexp.MatchString(constant.MatchIp, ip)
	return result
}
