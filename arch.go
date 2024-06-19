package cloudinfo

import "runtime"

func getArchFallback() string {
	v := runtime.GOARCH

	switch v {
	case "amd64":
		return "x86_64"
	default:
		return v
	}
}
