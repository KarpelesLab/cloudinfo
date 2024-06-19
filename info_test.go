package cloudinfo_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"testing"

	"github.com/KarpelesLab/cloudinfo"
)

func TestInfo(t *testing.T) {
	nfo, _ := cloudinfo.Load()
	log.Printf("info = %s", asJson(nfo))

	nfo = &cloudinfo.Info{
		Provider:     "test",
		AccountId:    "123456",
		Architecture: "x86_64",
		PublicIP:     []net.IP{net.ParseIP("1.2.3.4")},
		PrivateIP:    []net.IP{net.ParseIP("10.0.0.1")},
		Hostname:     "localhost.localdomain",
		Image:        "image-disk",
		ID:           "i-1232456",
		Type:         "testCloud",
		Location: cloudinfo.LocationArray{
			&cloudinfo.Location{"cloud", "test"},
			&cloudinfo.Location{"zone", "testzone"},
		},
		DMI: &cloudinfo.DMI{
			Cloud:       "test",
			Vendor:      "test123",
			ProductName: "unittest",
		},
	}
	log.Printf("info = %s", asJson(nfo))

	if nfo.Location.String() != "cloud=test,zone=testzone" {
		t.Errorf("invalid location string: %s", nfo.Location)
	}
}

func asJson(o any) []byte {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "    ")
	enc.Encode(o)
	return buf.Bytes()
}
