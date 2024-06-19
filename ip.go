package cloudinfo

import (
	"net"
	"net/netip"
)

type IPList []netip.Addr

func (l *IPList) addIP(ip net.IP) {
	a, ok := netip.AddrFromSlice(ip)
	if !ok {
		return
	}
	l.addAddr(a)
}

func (l *IPList) addAddr(a netip.Addr) {
	a = a.Unmap()

	for _, i := range *l {
		if i == a {
			return
		}
	}

	*l = append(*l, a)
}

func (l *IPList) addString(ip string) error {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return err
	}
	l.addAddr(a)
	return nil
}

// GetFirstV4 returns the first IPv4 found in the list
func (l IPList) GetFirstV4() (netip.Addr, bool) {
	for _, a := range l {
		if a.Is4() {
			return a, true
		}
	}
	return netip.Addr{}, false
}

// GetFirstV6 returns the first IPv6 found in the list
func (l IPList) GetFirstV6() (netip.Addr, bool) {
	for _, a := range l {
		if a.Is6() {
			return a, true
		}
	}
	return netip.Addr{}, false
}

// V4 returns an IPList with only IPv4 addresses included
func (l IPList) V4() IPList {
	var res IPList
	for _, a := range l {
		if a.Is4() {
			res = append(res, a)
		}
	}
	return res
}

// V6 returns an IPList with only IPv6 addresses included
func (l IPList) V6() IPList {
	var res IPList
	for _, a := range l {
		if a.Is6() {
			res = append(res, a)
		}
	}
	return res
}
