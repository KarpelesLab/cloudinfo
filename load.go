package cloudinfo

import (
	"errors"
	"fmt"
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
	info := sysinfo()
	cache := newCachedHttp()

	switch info.DMI.Cloud {
	case "aws":
		p := &awsProvider{info: info, cache: cache}
		return p.Fetch()
	case "gcp":
		p := &gcpProvider{cache: cache, info: info}
		return p.Fetch()
	case "scaleway":
		info.Type = info.DMI.ProductName
		info.fix()
		return info, errors.New("scaleway has no API for full machine information")
	default:
		info.fix()
		return info, fmt.Errorf("unsupported cloud provider %s", info.DMI.Cloud)
	}
}
