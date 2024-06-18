package cloudinfo

import "net"

type InfoLocation struct {
	Type  string `json:"type"`  // cloud, region, zone
	Value string `json:"value"` // the actual value
}

type Info struct {
	Provider     string          `json:"provider"`     // provider name code, such as "aws", "gcp", etc
	AccountId    string          `json:"account_id"`   // account id with provider
	Architecture string          `json:"architecture"` // x86_64
	PublicIP     []net.IP        `json:"public_ip"`
	PrivateIP    []net.IP        `json:"public_ip"`
	Hostname     string          `json:"hostname"`
	Image        string          `json:"image"`
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Location     []*InfoLocation `json:"location"` // structured location
}
