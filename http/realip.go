package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(http *http.Request) net.IP {
	// First try to get the IP from the X-Forwarded-For header if coming from a proxy
	ip := http.Header.Get("X-Forwarded-For")

	if ip != "" {
		parts := strings.SplitN(ip, ",", 2)
		part := strings.TrimSpace(parts[0])
		return net.ParseIP(part)
	}

	// Otherwise try to get the IP from the X-Real-IP if request coming from a load balancer
	ip = strings.TrimSpace(http.Header.Get("X-Real-IP"))

	if ip != "" {
		return net.ParseIP(ip)
	}

	// Otherwise, try to get the IP address from the RemoteAddr if coming from a physical devices
	address := strings.TrimSpace(http.RemoteAddr)
	host, _, err := net.SplitHostPort(address)

	if err != nil {
		return net.ParseIP(address)
	}

	// Otherwise, get the IP address from the host
	return net.ParseIP(host)
}
