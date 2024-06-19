package cloudinfo

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var (
	infoCache  *Info
	infoErr    error
	infoExpire time.Time
	infoLock   sync.RWMutex
)

// Load will load & return the info for the current machine. Even if an error happens, a Info structure will be
// returned containing some basic information.
func Load() (*Info, error) {
	if info, err := getInfoCache(); info != nil {
		return info, err
	}

	infoLock.Lock()
	defer infoLock.Unlock()

	if infoCache != nil && time.Until(infoExpire) >= 0 {
		return infoCache, infoErr
	}

	infoCache, infoErr = realLoad()

	if infoErr == nil {
		infoExpire = time.Now().Add(24 * time.Hour)
	} else {
		infoExpire = time.Now().Add(time.Hour)
	}

	return infoCache, infoErr
}

func getInfoCache() (*Info, error) {
	infoLock.RLock()
	defer infoLock.RUnlock()

	if time.Until(infoExpire) < 0 {
		return nil, nil
	}
	return infoCache, infoErr
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
