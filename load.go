package cloudinfo

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	infoCache *Info
	infoLock  sync.Mutex
)

// Load will load & return the info for the current machine. Even if an error happens, a Info structure will be
// returned containing some basic information.
// If no error is returned the info will be cached and the same Info object will be returned for each subsequent call
func Load() (*Info, error) {
	infoLock.Lock()
	defer infoLock.Unlock()

	if infoCache != nil {
		return infoCache, nil
	}

	info, err := realLoad()
	if err == nil {
		infoCache = info
	}
	return info, err
}

func realLoad() (*Info, error) {
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
