package cloudinfo

import (
	"fmt"
	"net"
)

type InfoLocation struct {
	Type  string `json:"type"`  // cloud, region, zone
	Value string `json:"value"` // the actual value
}

type Info struct {
	Provider     string          `json:"provider"`               // provider name code, such as "aws", "gcp", etc
	AccountId    string          `json:"account_id,omitempty"`   // account id with provider
	Architecture string          `json:"architecture,omitempty"` // x86_64
	PublicIP     []net.IP        `json:"public_ip,omitempty"`
	PrivateIP    []net.IP        `json:"public_ip,omitempty"`
	Hostname     string          `json:"hostname,omitempty"`
	Image        string          `json:"image,omitempty"`
	ID           string          `json:"id,omitempty"`
	Type         string          `json:"type,omitempty"`
	Location     []*InfoLocation `json:"location,omitempty"` // structured location
	DMI          *DMI            `json:"dmi,omitempty"`
}

// LoadInfo will load & return the info for the current machine. Even if an error happens, a Info structure will be
// returned containing some basic information.
func LoadInfo() (*Info, error) {
	dmi, _ := ReadDMI()
	info := &Info{Provider: dmi.Cloud, DMI: dmi}
	cache := newCachedHttp()

	switch dmi.Cloud {
	case "aws":
		p := &awsProvider{info: info, cache: cache}
		return p.Fetch()
	default:
		return info, fmt.Errorf("unsupported cloud provider %s", dmi.Cloud)
	}
}
