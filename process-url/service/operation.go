package service

import (
	"errors"
	"net/url"
	"strings"
)

type Operation int

const (
	OpUnknown Operation = iota
	OpRedirection
	OpCanonical
	OpAll
)

func (o Operation) String() string {
	switch o {
	case OpRedirection:
		return "redirection"
	case OpCanonical:
		return "canonical"
	case OpAll:
		return "all"
	default:
		return "unknown"
	}
}

func ParseOperation(s string) (Operation, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "redirection":
		return OpRedirection, nil
	case "canonical":
		return OpCanonical, nil
	case "all":
		return OpAll, nil
	default:
		return OpUnknown, errors.New("operation must be one of: redirection, canonical, all")
	}
}

func canonicalize(u *url.URL) {
	u.RawQuery = ""
	u.ForceQuery = false

	if p := u.EscapedPath(); p != "" && p != "/" {
		trimmed := strings.TrimRight(p, "/")
		u.Path = trimmed
	}
}

func redirectize(u *url.URL) {
	u.Host = "www.byfood.com"
}
