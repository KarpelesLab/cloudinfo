package cloudinfo

import (
	"os"
	"path/filepath"
	"strings"
)

// DMI can be used to find various information about the machine
type DMI struct {
	Cloud           string `json:"cloud"`
	Vendor          string `json:"sys_vendor"`
	Product         string `json:"product_name"`
	BoardAssetTag   string `json:"board_asset_tag"`
	ChassisAssetTag string `json:"chassis_asset_tag"`
}

func ReadDMI() (*DMI, error) {
	var err error
	var dmi *DMI

	dmi.Vendor, err = readDMI("sys_vendor")
	if err != nil {
		dmi.Cloud = "unknown"
		return dmi, err
	}

	switch dmi.Vendor {
	case "Amazon EC2":
		dmi.Cloud = "aws"
	case "Google":
		dmi.Cloud = "gcp"
	default:
		// Run this on the machine & send it as an issue to add support to that cloud provider
		// for foo in /sys/class/dmi/id/*; do echo $foo; cat $foo; done
		dmi.Cloud = "unknown"
	}

	dmi.BoardAssetTag, err = readDMI("board_asset_tag") // a uuid on google, the instance id on Amazon EC2
	if err != nil {
		return dmi, err
	}
	dmi.ChassisAssetTag, err = readDMI("chassis_asset_tag") // "Amazon EC2" on aws, nothing on google
	if err != nil {
		return dmi, err
	}

	// product_name on aws: m5a.large on google: Google Compute Engine
	dmi.Product, err = readDMI("product_name")
	if err != nil {
		return dmi, err
	}

	return dmi, nil
}

func readDMI(name string) (string, error) {
	res, err := os.ReadFile(filepath.Join("/sys/class/dmi/id", name))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(res)), nil
}
