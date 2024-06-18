package cloudinfo

import (
	"fmt"
	"net"
	"os"
)

// LoadInfo will load & return the info for the current machine. Even if an error happens, a Info structure will be
// returned containing some basic information.
func Load() (*Info, error) {
	dmi, _ := ReadDMI()
	info := &Info{
		Architecture: getArch(),
		Provider:     dmi.Cloud,
		DMI:          dmi,
	}
	if h, err := os.Hostname(); err == nil {
		info.Hostname = h
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
					info.addPrivateIP(a.IP)
					break
				}
				info.addPublicIP(a.IP)
			}
		}
	}
	cache := newCachedHttp()

	switch dmi.Cloud {
	case "aws":
		p := &awsProvider{info: info, cache: cache}
		return p.Fetch()
	case "gcp":
		p := &gcpProvider{cache: cache, info: info}
		return p.Fetch()
	default:
		return info, fmt.Errorf("unsupported cloud provider %s", dmi.Cloud)
	}
}
