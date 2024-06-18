package cloudinfo_test

import (
	"log"
	"testing"

	"github.com/KarpelesLab/cloudinfo"
)

func TestInfo(t *testing.T) {
	nfo, _ := cloudinfo.LoadInfo()

	log.Printf("info = %+v", nfo)
}
