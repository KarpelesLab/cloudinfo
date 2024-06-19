package cloudinfo

import (
	"net"
	"os"
	"strings"
)

type Info struct {
	Provider     string        `json:"provider"`               // provider name code, such as "aws", "gcp", etc
	AccountId    string        `json:"account_id,omitempty"`   // account id with provider
	Architecture string        `json:"architecture,omitempty"` // x86_64
	PublicIP     IPList        `json:"public_ip,omitempty"`
	PrivateIP    IPList        `json:"private_ip,omitempty"`
	Name         string        `json:"name,omitempty"`
	Hostname     string        `json:"hostname,omitempty"`
	Image        string        `json:"image,omitempty"`
	ID           string        `json:"id,omitempty"`
	Type         string        `json:"type,omitempty"`
	Location     LocationArray `json:"location,omitempty"` // structured location
	DMI          *DMI          `json:"dmi,omitempty"`
}

// Sysinfo loads and return information available from the local machine (including DMI)
func Sysinfo() *Info {
	info := sysinfo()
	info.fix()
	return info
}

func sysinfo() *Info {
	dmi, _ := ReadDMI()
	info := &Info{
		Architecture: getArch(),
		Provider:     dmi.Cloud,
		DMI:          dmi,
	}
	iflist, _ := net.Interfaces()
	for _, intf := range iflist {
		addrs, _ := intf.Addrs()
		for _, addr := range addrs {
			switch a := addr.(type) {
			case *net.IPNet:
				if a.IP.IsLoopback() {
					break
				}
				if a.IP.IsPrivate() || a.IP.IsLinkLocalUnicast() {
					info.PrivateIP.addIP(a.IP)
					break
				}
				info.PublicIP.addIP(a.IP)
			}
		}
	}

	return info
}

func (info *Info) fix() {
	if info.Hostname == "" {
		if h, err := os.Hostname(); err == nil {
			info.Hostname = h
		}
	}
	if info.Name == "" && info.Hostname != "" {
		str := info.Hostname
		if pos := strings.IndexByte(str, '.'); pos > 0 {
			str = str[:pos]
		}
		info.Name = str
	}
	if info.ID == "" {
		info.ID = info.DMI.ProductUUID
	}
}
