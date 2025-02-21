package cachecontrolheader_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mi-wada/cachecontrolheader"
)

func TestParse(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.header, func(t *testing.T) {
			t.Parallel()
			h := cachecontrolheader.Parse(tt.header)
			if diff := cmp.Diff(tt.want, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrict(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		header     string
		wantHeader *cachecontrolheader.Header
		wantErr    bool
	}{
		{
			header: "max-age=3600, must-revalidate, private",
			wantHeader: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header:  "max-age=3600, must-revalidate, private, unknown",
			wantErr: true,
		},
		{
			header:  "unknown",
			wantErr: true,
		},
		{
			header:  "unknown=10",
			wantErr: true,
		},
		{
			header:  "max-age=3600, must-revalidate, private, max-stale=invalid",
			wantErr: true,
		},
		{
			header:  "max-age=3600, must-revalidate, private, max-stale=10s",
			wantErr: true,
		},
		{
			header:  "max-age=invalid",
			wantErr: true,
		},
		{
			header:  "max-age=10s",
			wantErr: true,
		},
	} {
		tt := tt
		t.Run(tt.header, func(t *testing.T) {
			t.Parallel()
			h, err := cachecontrolheader.ParseStrict(tt.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error: %v, want: %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantHeader, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrict_IgnoreUnknownDirectives(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		header     string
		wantHeader *cachecontrolheader.Header
		wantErr    bool
	}{
		{
			header: "max-age=3600, must-revalidate, private, unknown",
			wantHeader: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header:     "unknown",
			wantHeader: &cachecontrolheader.Header{},
		},
		{
			header:     "unknown=10",
			wantHeader: &cachecontrolheader.Header{},
		},
		{
			header:  "max-age=invalid",
			wantErr: true,
		},
		{
			header:  "max-age=10s",
			wantErr: true,
		},
	} {
		tt := tt
		t.Run(tt.header, func(t *testing.T) {
			t.Parallel()
			h, err := cachecontrolheader.ParseStrict(tt.header, cachecontrolheader.IgnoreUnknownDirectives())
			if (err != nil) != tt.wantErr {
				t.Errorf("got error: %v, want: %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantHeader, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseStrict_IgnoreInvalidValues(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		header     string
		wantHeader *cachecontrolheader.Header
		wantErr    bool
	}{
		{
			header: "max-age=3600, must-revalidate, private, max-stale=invalid",
			wantHeader: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header: "max-age=3600, must-revalidate, private, max-stale=10s",
			wantHeader: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				MustRevalidate: true,
				Private:        true,
			},
		},
		{
			header:     "max-age=invalid",
			wantHeader: &cachecontrolheader.Header{},
		},
		{
			header:     "max-age=10s",
			wantHeader: &cachecontrolheader.Header{},
		},
		{
			header:  "unknown",
			wantErr: true,
		},
		{
			header:  "unknown=10",
			wantErr: true,
		},
	} {
		tt := tt
		t.Run(tt.header, func(t *testing.T) {
			t.Parallel()
			h, err := cachecontrolheader.ParseStrict(tt.header, cachecontrolheader.IgnoreInvalidValues())
			if (err != nil) != tt.wantErr {
				t.Errorf("got error: %v, want: %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.wantHeader, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_String(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.header.String(); got != tt.want {
				t.Errorf("Header.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
