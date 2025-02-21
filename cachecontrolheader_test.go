package cachecontrolheader_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mi-wada/cachecontrolheader"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		header string
		want   *cachecontrolheader.Header
	}{
		{
			header: "max-age=3600, must-revalidate, private",
			want: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header: "",
			want:   &cachecontrolheader.Header{},
		},
		{
			header: "unknown",
			want:   &cachecontrolheader.Header{},
		},
		{
			header: "unknown=10",
			want:   &cachecontrolheader.Header{},
		},
		{
			header: "max-age=invalid",
			want:   &cachecontrolheader.Header{},
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h := cachecontrolheader.Parse(tt.header)
			if diff := cmp.Diff(tt.want, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrict(t *testing.T) {
	for _, tt := range []struct {
		header string
	}{
		{
			header: "max-age=3600, must-revalidate, private, unknown",
		},
		{
			header: "unknown",
		},
		{
			header: "unknown=10",
		},
		{
			header: "max-age=invalid, must-revalidate, private",
		},
		{
			header: "max-age=10s, must-revalidate, private",
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := cachecontrolheader.ParseStrict(tt.header)
			if err == nil {
				t.Errorf("want error, but got nil. Header struct: %v", h)
			}
		})
	}
}

func TestParseStrict_IgnoreUnknownDirectives(t *testing.T) {
	for _, tt := range []struct {
		header string
		want   *cachecontrolheader.Header
	}{
		{
			header: "max-age=3600, must-revalidate, private, unknown",
			want: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header: "unknown",
			want:   &cachecontrolheader.Header{},
		},
		{
			header: "unknown=10",
			want:   &cachecontrolheader.Header{},
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := cachecontrolheader.ParseStrict(tt.header, cachecontrolheader.IgnoreUnknownDirectives())
			if err != nil {
				t.Errorf("want nil, but got error: %v", err)
			}
			if diff := cmp.Diff(tt.want, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrict_IgnoreInvalidValues(t *testing.T) {
	for _, tt := range []struct {
		header string
		want   *cachecontrolheader.Header
	}{
		{
			header: "max-age=invalid, must-revalidate, private",
			want: &cachecontrolheader.Header{
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header: "max-age=10s, must-revalidate, private",
			want: &cachecontrolheader.Header{
				MustRevalidate: true,
				Private:        true,
			},
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := cachecontrolheader.ParseStrict(tt.header, cachecontrolheader.IgnoreInvalidValues())
			if err != nil {
				t.Errorf("want nil, but got error: %v", err)
			}
			if diff := cmp.Diff(tt.want, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_String(t *testing.T) {
	for _, tt := range []struct {
		header *cachecontrolheader.Header
		want   string
	}{
		{
			header: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
			want: "max-age=3600, must-revalidate, private",
		},
		{
			header: &cachecontrolheader.Header{},
			want:   "",
		},
	} {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.header.String(); got != tt.want {
				t.Errorf("Header.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
