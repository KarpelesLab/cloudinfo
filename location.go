package cloudinfo

import (
	"fmt"
	"strings"
)

type Location struct {
	Type  string `json:"type"`  // cloud, region, zone
	Value string `json:"value"` // the actual value
}

type LocationArray []*Location

func makeLocation(s ...string) LocationArray {
	var l LocationArray

	ln := len(s)

	for i := 0; i < ln; i += 2 {
		l = append(l, &Location{Type: s[i], Value: s[i+1]})
	}
	return l
}

func (l *Location) String() string {
	return fmt.Sprintf("%s=%s", l.Type, l.Value)
}

func (la LocationArray) String() string {
	s := make([]string, len(la))
	for n, v := range la {
		s[n] = v.String()
	}
	return strings.Join(s, ",")
}
