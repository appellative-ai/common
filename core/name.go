package core

import (
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"
)

const (
	Fragment = "#"
	Colon    = ":"
	Slash    = "/"
)

var (
	counter = new(atomic.Int64)
)

type Name struct {
	Collective string `json:"collective"`
	Domain     string `json:"domain"`
	Kind       string `json:"kind"`
	Path       string `json:"path"`
	Fragment   string `json:"fragment"`
}

func NewName(name string) Name {
	return parse(name)
}

func AddFragment(name, fragment string) string {
	return addFragment(name, fragment)
}

func Versioned(name string) string {
	return AddFragment(name, fmt.Sprintf("%v", counter.Add(1)))
}

/*
	func Collective(name string) string {
		return NewName(name).Collective
	}

	func Domain(name string) string {
		return NewName(name).Domain
	}


	func Path(name string) string {
		return NewName(name).Path
	}

	func Fragment(name string) string {
		return NewName(name).Path
	}
*/

func Kind(name string) string {
	return NewName(name).Kind
}

func parse(name string) Name {
	if name == "" {
		return Name{}
	}
	u, err := url.Parse(name)
	if err != nil {
		return Name{Collective: err.Error()}
	}
	n := Name{Collective: u.Scheme, Fragment: u.Fragment}
	i := strings.Index(u.Opaque, Colon)
	if i < 0 {
		n.Collective = "error, missing second colon"
		return n
	}
	n.Domain = u.Opaque[:i]
	path := u.Opaque[i:]
	i = strings.Index(path, Slash)
	if i < 0 {
		return n
	}
	n.Kind = path[1:i]
	n.Path = path[i:]
	return n
}

func addFragment(name, fragment string) string {
	if name == "" {
		return name
	}
	i := strings.Index(name, Fragment)
	if i == -1 {
		return name + Fragment + fragment
	}
	return name
}
