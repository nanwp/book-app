package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

func ProcessURL(ctx context.Context, rawURL string, operation string) (string, error) {
	log := log.Ctx(ctx).With().Str("service", "process_url").Logger()

	u, err := url.Parse(rawURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse URL")
		return "", errors.New("invalid url")
	}
	if u.Scheme == "" || (u.Scheme != "http" && u.Scheme != "https") {
		log.Error().Err(err).Msg("URL must include scheme http or https")
		return "", errors.New("url must include scheme http or https")
	}
	if u.Host == "" {
		log.Error().Msg("URL must include a host")
		return "", errors.New("url must include a host")
	}

	op, err := ParseOperation(operation)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse operation")
		return "", err
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
