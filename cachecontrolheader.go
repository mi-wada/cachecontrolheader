package cachecontrolheader

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Header represents a Cache-Control header.
type Header struct {
	NoCache         bool
	NoStore         bool
	NoTransform     bool
	OnlyIfCached    bool
	MustRevalidate  bool
	MustUnderstand  bool
	Private         bool
	ProxyRevalidate bool
	Public          bool
	// In request header, it means <https://datatracker.ietf.org/doc/html/rfc9111.html#section-5.2.1.1>.
	// In response header, it means <>
	MaxAge   time.Duration
	MaxStale time.Duration
	MinFresh time.Duration
	SMaxAge  time.Duration
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
