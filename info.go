package cloudinfo

type Info struct {
	Provider     string        `json:"provider"`               // provider name code, such as "aws", "gcp", etc
	AccountId    string        `json:"account_id,omitempty"`   // account id with provider
	Architecture string        `json:"architecture,omitempty"` // x86_64
	PublicIP     IPList        `json:"public_ip,omitempty"`
	PrivateIP    IPList        `json:"private_ip,omitempty"`
	Hostname     string        `json:"hostname,omitempty"`
	Image        string        `json:"image,omitempty"`
	ID           string        `json:"id,omitempty"`
	Type         string        `json:"type,omitempty"`
	Location     LocationArray `json:"location,omitempty"` // structured location
	DMI          *DMI          `json:"dmi,omitempty"`
}
