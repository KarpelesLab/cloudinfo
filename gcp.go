package cloudinfo

import (
	"encoding/json"
	"net"
	"path"
	"strconv"
	"strings"
)

type gcpProvider struct {
	cache *cachedHttp
	info  *Info
}

type gcpMetadata struct {
	// attributes: physical_host, ssh-keys, stop-state
	CPUPlatform string `json:"cpuPlatform"` // Intel Broadwell
	Description string `json:"description"` // ?
	Hostname    string `json:"hostname"`    // instance-name.asia-northeast1-b.c.project-name.internal
	ID          int64  `json:"id"`          // 123456...
	Image       string `json:"image"`       // projects/debian-cloud/global/images/debian-12-bookworm-v20240515
	MachineType string `json:"machineType"` // projects/<id>/machineTypes/e2-standard-4
	Name        string `json:"name"`        // name of machine
	// "licenses": [{"id": "123456..."}]
	NetworkInterfaces []struct {
		AccessConfigs []struct {
			ExternalIP string `json:"externalIp"`
			Type       string `json:"type"` // ONE_TO_ONE_NAT
		} `json:"accessConfigs"`
		// "dnsServers": ["169.254.169.254"]
		// "forwardedIps": [],
		Gateway string `json:"gateway"`
		IP      string `json:"ip"` // local ip (private)
		// "ipAliases": [],
		// "mac": "42:01:0a:92:00:02",
		// "mtu": 1460,
		// "network": "projects/<id>/networks/default",
		// "subnetmask": "255.255.240.0",
		// "targetInstanceIps": []
	} `json:"networkInterfaces"`
	//Tags         []string `json:"tags"`
	Zone string `json:"zone"` // projects/<id>/zones/asia-northeast1-b
}

func (g *gcpProvider) Name() string {
	return "gcp"
}

func (g *gcpProvider) Fetch() (*Info, error) {
	res, _, err := g.cache.GetWithHeaders("http://metadata.google.internal/computeMetadata/v1/instance/?recursive=true&timeout_sec=1", map[string]string{"Metadata-Flavor": "Google"})
	if err != nil {
		return g.info, err
	}
	var info *gcpMetadata
	err = json.Unmarshal(res, &info)
	if err != nil {
		return g.info, err
	}

	g.info.Hostname = info.Hostname
	g.info.ID = strconv.FormatInt(info.ID, 10)
	g.info.Image = info.Image
	g.info.Type = path.Base(info.MachineType)
	if a := strings.Split(info.Zone, "/"); len(a) >= 2 && a[0] == "projects" {
		g.info.AccountId = a[1]
	}

	for _, intf := range info.NetworkInterfaces {
		if ip := net.ParseIP(intf.IP); ip != nil {
			g.info.addPrivateIP(ip)
		}
		for _, ac := range intf.AccessConfigs {
			if ip := net.ParseIP(ac.ExternalIP); ip != nil {
				g.info.addPublicIP(ip)
			}
		}
	}

	zone := path.Base(info.Zone)
	region := zone
	if pos := strings.LastIndexByte(region, '-'); pos > 0 {
		region = region[:pos]
	}

	g.info.Location = makeLocation("cloud", "gcp", "region", region, "zone", zone)

	return g.info, nil
}
