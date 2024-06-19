package cloudinfo

import (
	"os"
	"path/filepath"
	"strings"
)

// DMI can be used to find various information about the machine
type DMI struct {
	Cloud           string `json:"cloud"`
	Vendor          string `json:"sys_vendor,omitempty"`
	ProductName     string `json:"product_name,omitempty"`
	ProductUUID     string `json:"product_uuid,omitempty"`
	ProductVersion  string `json:"product_version,omitempty"`
	BoardAssetTag   string `json:"board_asset_tag,omitempty"`
	ChassisAssetTag string `json:"chassis_asset_tag,omitempty"`
}

func ReadDMI() (*DMI, error) {
	dmi := &DMI{Cloud: "unknown"}

	dmi.Vendor = readDMI("sys_vendor")
	dmi.BoardAssetTag = readDMI("board_asset_tag")     // a uuid on google, the instance id on Amazon EC2
	dmi.ChassisAssetTag = readDMI("chassis_asset_tag") // "Amazon EC2" on aws, nothing on google
	// product_name on aws: m5a.large on google: Google Compute Engine
	dmi.ProductName = readDMI("product_name")
	dmi.ProductUUID = readDMI("product_uuid")
	dmi.ProductVersion = readDMI("product_version")

	switch strings.ToLower(dmi.Vendor) {
	case "amazon ec2":
		dmi.Cloud = "aws"
	case "google":
		dmi.Cloud = "gcp"
	case "scaleway":
		dmi.Cloud = "scaleway"
	case "xen":
		// ProductVersion = 4.11.amazon
		if strings.Contains(dmi.ProductVersion, "amazon") {
			dmi.Cloud = "aws"
			break
		}
	}
	// if Cloud=unknown, run this on the machine & send it as an issue to add support to that cloud provider
	// for foo in /sys/class/dmi/id/*; do echo $foo; cat $foo; done

	return dmi, nil
}

func readDMI(name string) string {
	res, err := os.ReadFile(filepath.Join("/sys/class/dmi/id", name))
	if err != nil {
		return ""
	}
	resStr := strings.TrimSpace(string(res))
	switch resStr {
	case "To be filled by O.E.M.", "Default string":
		return ""
	default:
		return resStr
	}
}
