package cachecontrolheader

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Header represents a Cache-Control header.
type Header struct {
	// In request header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.4
	// In response header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.4
	NoCache bool
	// In request header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.5
	// In response header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.5
	NoStore bool
	// In request header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.6
	// In response header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.6
	NoTransform bool
	// In request header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.7
	OnlyIfCached bool
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.2
	MustRevalidate bool
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.3
	MustUnderstand bool
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.7
	Private bool
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.8
	ProxyRevalidate bool
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.9
	Public bool
	// In request header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.1
	// In response header: https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.1
	MaxAge time.Duration
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.2
	MaxStale time.Duration
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.3
	MinFresh time.Duration
	// https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.2.10
	SMaxAge time.Duration
}

// Parse parses a Cache-Control header based on [RFC9111](https://datatracker.ietf.org/doc/html/rfc9111.html).
func Parse(header string) (*Header, error) {
	header = strings.ToLower(strings.ReplaceAll(header, " ", ""))

	h := Header{}
	if header == "" {
		return &h, nil
	}
	directives := strings.Split(header, ",")
	for _, d := range directives {
		splited := strings.SplitN(d, "=", 2)
		switch len(splited) {
		case 1:
			switch splited[0] {
			case "no-cache":
				h.NoCache = true
			case "no-store":
				h.NoStore = true
			case "only-if-cached":
				h.OnlyIfCached = true
			case "must-revalidate":
				h.MustRevalidate = true
			case "must-understand":
				h.MustUnderstand = true
			case "private":
				h.Private = true
			case "proxy-revalidate":
				h.ProxyRevalidate = true
			case "public":
				h.Public = true
			default:
				return nil, fmt.Errorf("unknown directive: %s", splited[0])
			}
		case 2:
			k := splited[0]
			v, err := time.ParseDuration(strings.TrimSpace(splited[1]) + "s")
			if err != nil {
				return nil, fmt.Errorf("failed to parse duration for directive %s=%s: %w", splited[0], splited[1], err)
			}
			switch k {
			case "max-age":
				h.MaxAge = v
			case "max-stale":
				h.MaxStale = v
			case "min-fresh":
				h.MinFresh = v
			case "s-maxage":
				h.SMaxAge = v
			default:
				return nil, fmt.Errorf("unknown directive: %s", k)
			}
		}
	}
	return &h, nil
}

// ParseReader parses a Cache-Control header based on [RFC9111](https://datatracker.ietf.org/doc/html/rfc9111.html) from an [io.Reader].
func ParseReader(r io.Reader) (*Header, error) {
	header, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(string(header))
}
