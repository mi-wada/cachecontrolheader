// Package cachecontrolheader provides functionality to parse and handle
// Cache-Control headers based on RFC 9111 Section 5.2.
package cachecontrolheader

import (
	"fmt"
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

// Parse parses a Cache-Control header based on RFC 9111 Section 5.2.
// By default, it ignores unknown directives and invalid values.
// To return an error when those cases, use [ParseStrict] instead.
func Parse(header string) *Header {
	h, _ := parse(header, IgnoreInvalidValues(), IgnoreUnknownDirectives())
	return h
}

// ParseStrict strictly parses a Cache-Control header based on RFC 9111 Section 5.2.
// By default, it returns an error when unknown directives or invalid values are found.
// To ignore either of them, use [IgnoreUnknownDirectives] or [IgnoreInvalidValues] options.
// To ignore both, use [Parse] instead.
func ParseStrict(header string, opts ...parseOption) (*Header, error) {
	return parse(header, opts...)
}

// IgnoreUnknownDirectives allows to ignore unknown directives.
func IgnoreUnknownDirectives() parseOption {
	return func(o *option) {
		o.ignoreUnknownDirectives = true
	}
}

// IgnoreInvalidValues allows to ignore directives that have invalid values.
// Invalid values examples: `max-age=invalid`, `max-stale=1s`
func IgnoreInvalidValues() parseOption {
	return func(o *option) {
		o.ignoreInvalidValues = true
	}
}

type option struct {
	ignoreUnknownDirectives bool
	ignoreInvalidValues     bool
}
type parseOption func(*option)

// Header represents a Cache-Control header.
type Header struct {
	MaxAge          *time.Duration // max-age directive
	MaxStale        *time.Duration // max-stale directive
	MinFresh        *time.Duration // min-fresh directive
	NoCache         bool           // no-cache directive
	NoStore         bool           // no-store directive
	NoTransform     bool           // no-transform directive
	OnlyIfCached    bool           // only-if-cached directive
	MustRevalidate  bool           // must-revalidate directive
	MustUnderstand  bool           // must-understand directive
	Private         bool           // private directive
	ProxyRevalidate bool           // proxy-revalidate directive
	Public          bool           // public directive
	SMaxAge         *time.Duration // s-maxage directive
}

// String returns a string representation of the Cache-Control header.
func (h *Header) String() string {
	var ds []string
	if h.MaxAge != nil {
		ds = append(ds, fmt.Sprintf("%s=%d", dMaxAge, int(h.MaxAge.Seconds())))
	}
	if h.MaxStale != nil {
		ds = append(ds, fmt.Sprintf("%s=%d", dMaxStale, int(h.MaxStale.Seconds())))
	}
	if h.MinFresh != nil {
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
	if h.SMaxAge != nil {
		ds = append(ds, fmt.Sprintf("%s=%d", dSMaxAge, int(h.SMaxAge.Seconds())))
	}
	return strings.Join(ds, ", ")
}

// parse parses a Cache-Control header based on RFC 9111 Section 5.2.
// By default, it returns an error when unknown directives found.
// To ignore unknown directives, use [IgnoreUnknownDirectives] option.
// By default, it returns an error when invalid values found.
// To ignore invalid values, use [IgnoreInvalidValues] option.
func parse(header string, opts ...parseOption) (*Header, error) {
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
				if option.ignoreUnknownDirectives {
					continue
				}
				return nil, fmt.Errorf("unknown directive: %s", splited[0])
			}
		case 2:
			k := splited[0]
			v, err := time.ParseDuration(strings.TrimSpace(splited[1]) + "s")
			if err != nil {
				if option.ignoreInvalidValues {
					continue
				} else {
					return nil, fmt.Errorf("failed to parse the value of directive(%s=%s): %w", splited[0], splited[1], err)
				}
			}
			switch k {
			case dMaxAge:
				h.MaxAge = &v
			case dMaxStale:
				h.MaxStale = &v
			case dMinFresh:
				h.MinFresh = &v
			case dSMaxAge:
				h.SMaxAge = &v
			default:
				if option.ignoreUnknownDirectives {
					continue
				}
				return nil, fmt.Errorf("unknown directive: %s", k)
			}
		}
	}
	return &h, nil
}
