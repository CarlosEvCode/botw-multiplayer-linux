package network

import (
	"net"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

type NetworkInfo struct {
	ZeroTierIP  string
	TailscaleIP string
	LocalIP     string
	PublicIP    string
}

func GetNetworkInfo() NetworkInfo {
	info := NetworkInfo{
		ZeroTierIP:  getZeroTierIP(),
		TailscaleIP: getTailscaleIP(),
		LocalIP:     getLocalIP(),
	}
	return info
}

func getZeroTierIP() string {
	// Try zerotier-cli listnetworks first
	out, err := exec.Command("zerotier-cli", "listnetworks").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			// Format: 200 listnetworks <nwid> <name> <mac> <status> <type> <dev> <assigned-ips>
			if len(fields) >= 8 && fields[4] == "OK" {
				for _, f := range fields[7:] {
					ips := strings.Split(f, ",")
					for _, ipMask := range ips {
						ip := strings.Split(ipMask, "/")[0]
						parsed := net.ParseIP(ip)
						if parsed != nil && parsed.To4() != nil {
							return ip
						}
					}
				}
			}
		}
	}

	// Fallback to checking any zt* or zerotier* network interface
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if strings.HasPrefix(iface.Name, "zt") || strings.HasPrefix(iface.Name, "zerotier") {
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
	}
	return ""
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

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if !strings.HasPrefix(ip, "100.") { // Exclude Tailscale CGNAT 100.x
					return ip
				}
			}
		}
	}
	return "127.0.0.1"
}

func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

