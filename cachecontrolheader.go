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

// directives
const (
	dMaxAge          = "max-age"
	dMaxStale        = "max-stale"
	dMinFresh        = "min-fresh"
	dNoCache         = "no-cache"
	dNoStore         = "no-store"
	dNoTransform     = "no-transform"
	dOnlyIfCached    = "only-if-cached"
	dMustRevalidate  = "must-revalidate"
	dMustUnderstand  = "must-understand"
	dPrivate         = "private"
	dProxyRevalidate = "proxy-revalidate"
	dPublic          = "public"
	dSMaxAge         = "s-maxage"
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

// String returns a string representation of the Cache-Control header.
func (h *Header) String() string {
	var ds []string
	if h.MaxAge > 0 {
		ds = append(ds, fmt.Sprintf("%s=%d", dMaxAge, int(h.MaxAge.Seconds())))
	}
	if h.MaxStale > 0 {
		ds = append(ds, fmt.Sprintf("%s=%d", dMaxStale, int(h.MaxStale.Seconds())))
	}
	if h.MinFresh > 0 {
		ds = append(ds, fmt.Sprintf("%s=%d", dMinFresh, int(h.MinFresh.Seconds())))
	}
	if h.NoCache {
		ds = append(ds, dNoCache)
	}
	if h.NoStore {
		ds = append(ds, dNoStore)
	}
	if h.NoTransform {
		ds = append(ds, dNoTransform)
	}
	if h.OnlyIfCached {
		ds = append(ds, dOnlyIfCached)
	}
	if h.MustRevalidate {
		ds = append(ds, dMustRevalidate)
	}
	if h.MustUnderstand {
		ds = append(ds, dMustUnderstand)
	}
	if h.Private {
		ds = append(ds, dPrivate)
	}
	if h.ProxyRevalidate {
		ds = append(ds, dProxyRevalidate)
	}
	if h.Public {
		ds = append(ds, dPublic)
	}
	if h.SMaxAge > 0 {
		ds = append(ds, fmt.Sprintf("%s=%d", dSMaxAge, int(h.SMaxAge.Seconds())))
	}
	return strings.Join(ds, ", ")
}

type option struct {
	errorOnUnknownDirectives bool
	errorOnInvalidValues     bool
}
type parseOption func(*option)

// ErrorOnUnknown allows to return an error when unknown directives are found.
func ErrorOnUnknown() parseOption {
	return func(o *option) {
		o.errorOnUnknownDirectives = true
	}
}

// ErrorOnInvalidValues allows to return an error when invalid values are found.
// Invalid values examples: `max-age=invalid`, `max-stale=1s`
func ErrorOnInvalidValues() parseOption {
	return func(o *option) {
		o.errorOnInvalidValues = true
	}
}

// Parse parses a Cache-Control header based on RFC9111.
// By default, it ignores unknown directives.
// To return an error when unknown directives are found, use [ErrorOnUnknown] option.
// By default, it ignores directives that have invalid values, like `max-age=invalid`.
// To return an error when invalid values are found, use [ErrorOnInvalidValues] option.
func Parse(header string, opts ...parseOption) (*Header, error) {
	option := option{}
	for _, opt := range opts {
		opt(&option)
	}
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
			case dNoCache:
				h.NoCache = true
			case dNoStore:
				h.NoStore = true
			case dOnlyIfCached:
				h.OnlyIfCached = true
			case dMustRevalidate:
				h.MustRevalidate = true
			case dMustUnderstand:
				h.MustUnderstand = true
			case dPrivate:
				h.Private = true
			case dProxyRevalidate:
				h.ProxyRevalidate = true
			case dPublic:
				h.Public = true
			default:
				if option.errorOnUnknownDirectives {
					return nil, fmt.Errorf("unknown directive: %s", splited[0])
				}
			}
		case 2:
			k := splited[0]
			v, err := time.ParseDuration(strings.TrimSpace(splited[1]) + "s")
			if err != nil && option.errorOnInvalidValues {
				return nil, fmt.Errorf("failed to parse duration for directive %s=%s: %w", splited[0], splited[1], err)
			}
			switch k {
			case dMaxAge:
				h.MaxAge = v
			case dMaxStale:
				h.MaxStale = v
			case dMinFresh:
				h.MinFresh = v
			case dSMaxAge:
				h.SMaxAge = v
			default:
				if option.errorOnUnknownDirectives {
					return nil, fmt.Errorf("unknown directive: %s", k)
				}
			}
		}
	}
	return &h, nil
}

// ParseReader parses a Cache-Control header from an io.Reader based on RFC9111.
// By default, it ignores unknown directives.
// To return an error when unknown directives are found, use [ErrorOnUnknown] option.
func ParseReader(r io.Reader, opts ...parseOption) (*Header, error) {
	header, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return Parse(string(header), opts...)
}
