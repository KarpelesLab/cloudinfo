//go:build !linux

package cloudinfo

func getArch() string {
	return getArchFallback()
}
