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
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := cachecontrolheader.Parse(tt.header)
			if err != nil {
				t.Fatal(err)
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
