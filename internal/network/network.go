package network

import (
	"net"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

type NetworkInfo struct {
	LocalIP     string
	TailscaleIP string
	ZeroTierIP  string
	PublicIP    string
}

func GetNetworkInfo() NetworkInfo {
	info := NetworkInfo{
		LocalIP:     getLocalIP(),
		TailscaleIP: getTailscaleIP(),
		ZeroTierIP:  getZeroTierIP(),
	}
	return info
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if !strings.HasPrefix(ip, "100.") { // Tailscale typically uses 100.x
					return ip
				}
			}
		}
	}
	return "127.0.0.1"
}

func getTailscaleIP() string {
	// Try tailscale CLI first
	out, err := exec.Command("tailscale", "ip", "-4").Output()
	if err == nil {
		ip := strings.TrimSpace(string(out))
		if ip != "" {
			return ip
		}
	}

	// Fallback to checking tailscale0 interface
	iface, err := net.InterfaceByName("tailscale0")
	if err == nil {
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

func getZeroTierIP() string {
	// Check for any zt* interface
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "zt") {
			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}
