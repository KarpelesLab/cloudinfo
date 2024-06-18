package cloudinfo

import (
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
	PrivateIP    []net.IP        `json:"private_ip,omitempty"`
	Hostname     string          `json:"hostname,omitempty"`
	Image        string          `json:"image,omitempty"`
	ID           string          `json:"id,omitempty"`
	Type         string          `json:"type,omitempty"`
	Location     []*InfoLocation `json:"location,omitempty"` // structured location
	DMI          *DMI            `json:"dmi,omitempty"`
}

func (i *Info) addPrivateIP(ip net.IP) {
	for _, prev := range i.PrivateIP {
		if prev.Equal(ip) {
			return
		}
	}
	i.PrivateIP = append(i.PrivateIP, ip)
}

func (i *Info) addPublicIP(ip net.IP) {
	for _, prev := range i.PublicIP {
		if prev.Equal(ip) {
			return
		}
	}
	i.PublicIP = append(i.PublicIP, ip)
}
