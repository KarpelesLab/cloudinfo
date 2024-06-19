//go:build linux

package cloudinfo

import "syscall"

func getArch() string {
	var name syscall.Utsname
	if err := syscall.Uname(&name); err != nil {
		return getArchFallback()
	}

	return fncUnameHelperToString(name.Machine)
}

// fncUnameHelperToString converts a uname string to a go string
func fncUnameHelperToString[T int8 | uint8](v [65]T) string {
	out := make([]byte, len(v))
	for i := 0; i < len(v); i++ {
		if v[i] == 0 {
			return string(out[:i])
		}
		out[i] = byte(v[i])
	}
	return string(out)
}
