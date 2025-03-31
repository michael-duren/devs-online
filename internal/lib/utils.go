package lib

import (
	"net"
)

func IsValidIP(ip string) (bool, error) {
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		host = ip
	}

	parsedIp := net.ParseIP(host)
	if parsedIp == nil {
		return false, net.InvalidAddrError("invalid ip address")
	}
	if parsedIp.To4() != nil {
		return true, nil
	}
	return false, net.InvalidAddrError("valid ipv6 address but only ipv4 is supported")
}
