// Package cachecontrolheader provides functionality to parse and handle
// Cache-Control headers based on RFC9111.
// https://datatracker.ietf.org/doc/html/rfc9111.html#name-cache-control
package cachecontrolheader

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Header represents a Cache-Control header.
type Header struct {
	MaxAge          time.Duration // max-age directive
	MaxStale        time.Duration // max-stale directive
	MinFresh        time.Duration // min-fresh directive
	NoCache         bool          // no-cache directive
	NoStore         bool          // no-store directive
	NoTransform     bool          // no-transform directive
	OnlyIfCached    bool          // only-if-cached directive
	MustRevalidate  bool          // must-revalidate directive
	MustUnderstand  bool          // must-understand directive
	Private         bool          // private directive
	ProxyRevalidate bool          // proxy-revalidate directive
	Public          bool          // public directive
	SMaxAge         time.Duration // s-maxage directive
}

func (h *Header) String() string {
	var ds []string
	if h.MaxAge > 0 {
		ds = append(ds, fmt.Sprintf("max-age=%d", int(h.MaxAge.Seconds())))
	}
	if h.MaxStale > 0 {
		ds = append(ds, fmt.Sprintf("max-stale=%d", int(h.MaxStale.Seconds())))
	}
	if h.MinFresh > 0 {
		ds = append(ds, fmt.Sprintf("min-fresh=%d", int(h.MinFresh.Seconds())))
	}
	if h.NoCache {
		ds = append(ds, "no-cache")
	}
	if h.NoStore {
		ds = append(ds, "no-store")
	}
	if h.NoTransform {
		ds = append(ds, "no-transform")
	}
	if h.OnlyIfCached {
		ds = append(ds, "only-if-cached")
	}
	if h.MustRevalidate {
		ds = append(ds, "must-revalidate")
	}
	if h.MustUnderstand {
		ds = append(ds, "must-understand")
	}
	if h.Private {
		ds = append(ds, "private")
	}
	if h.ProxyRevalidate {
		ds = append(ds, "proxy-revalidate")
	}
	if h.Public {
		ds = append(ds, "public")
	}
	if h.SMaxAge > 0 {
		ds = append(ds, fmt.Sprintf("s-maxage=%d", int(h.SMaxAge.Seconds())))
	}
	return strings.Join(ds, ", ")
}

// Parse parses a Cache-Control header based on RFC9111.
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

// ParseReader parses a Cache-Control header from an io.Reader based on RFC9111.
func ParseReader(r io.Reader) (*Header, error) {
	header, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(string(header))
}
