package cloudinfo

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
)

type awsProvider struct {
	cache *cachedHttp
	info  *Info
	token string
	err   error
}

type awsIdentity struct {
	AccountId        string `json:"accountId"`        // 12 digits
	Architecture     string `json:"architecture"`     // x86_64
	AvailabilityZone string `json:"availabilityZone"` // ap-northeast-1c
	// "billingProducts" : null,
	// "devpayProductCodes" : null,
	// "marketplaceProductCodes" : null,
	ImageId      string `json:"imageId"`      // ami-xxx
	InstanceId   string `json:"instanceId"`   // i-xxx
	InstanceType string `json:"instanceType"` // m5a.large
	// "kernelId" : null,
	// "pendingTime" : "2021-10-17T13:39:07Z",
	PrivateIp string `json:"privateIp"`
	// "ramdiskId" : null,
	Region  string `json:"region"`  // ap-northeast-1
	Version string `json:"version"` // 2017-09-30
}

func (a *awsProvider) Name() string {
	return "aws"
}

func (a *awsProvider) Fetch() error {
	// initialize a.res
	if a.info == nil {
		a.info = &Info{}
	}

	err := a.getToken()
	if err != nil {
		return err
	}
	err = a.getIdentity()
	if err != nil {
		return err
	}

	return nil
}

// getToken fetches an api token from the aws server and stores it into a.token
func (a *awsProvider) getToken() error {
	if a.token != "" {
		// already have a token
		return nil
	}

	tokenB, _, err := a.cache.PutWithHeaders("http://169.254.169.254/latest/api/token", map[string]string{"X-aws-ec2-metadata-token-ttl-seconds": "60"})
	if err != nil {
		return err
	}
	token := strings.TrimSpace(string(tokenB))
	if token == "" {
		return errors.New("could not fetch aws token")
	}
	a.token = token
	return nil
}

func (a *awsProvider) getMeta(p string) (string, error) {
	res, _, err := a.cache.GetWithHeaders("http://169.254.169.254/latest/meta-data/"+p, map[string]string{"X-aws-ec2-metadata-token": a.token})
	return string(res), err
}

func (a *awsProvider) getIdentity() error {
	res, _, err := a.cache.GetWithHeaders("http://169.254.169.254/latest/dynamic/instance-identity/document", map[string]string{"X-aws-ec2-metadata-token": a.token})
	if err != nil {
		return err
	}
	var info *awsIdentity
	err = json.Unmarshal(res, &info)
	if err != nil {
		return err
	}

	// fill a.info with the info here
	a.info.AccountId = info.AccountId
	a.info.Architecture = info.Architecture
	a.info.Image = info.ImageId
	a.info.ID = info.InstanceId
	a.info.Type = info.InstanceType
	if ip := net.ParseIP(info.PrivateIp); ip != nil {
		a.info.PrivateIP = append(a.info.PrivateIP, ip)
	}

	a.info.Location = []*InfoLocation{
		&InfoLocation{Type: "cloud", Value: "aws"},
		&InfoLocation{Type: "region", Value: info.Region},
		&InfoLocation{Type: "zone", Value: info.AvailabilityZone},
	}

	return nil
}
