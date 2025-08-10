package service

import (
	"errors"
	"net/url"
	"strings"
)

func ProcessURL(rawURL string, op Operation) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.New("invalid url")
	}
	if u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return "", errors.New("url must include scheme http or https")
	}
	if u.Host == "" {
		return "", errors.New("url must include a host")
	}

	switch op {
	case OpCanonical:
		canonicalize(u)
		return u.String(), nil
	case OpRedirection:
		redirectize(u)
		return strings.ToLower(u.String()), nil
	case OpAll:
		canonicalize(u)
		redirectize(u)
		return strings.ToLower(u.String()), nil
	default:
		return "", errors.New("unsupported operation")
	}
}
