package cachecontrolheader_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/mi-wada/cachecontrolheader"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		header     string
		wantHeader *cachecontrolheader.Header
	}{
		{
			header: "max-age=3600, private, must-revalidate",
			wantHeader: &cachecontrolheader.Header{
				MaxAge:         3600 * time.Second,
				Private:        true,
				MustRevalidate: true,
			},
		},
		{
			header:     "",
			wantHeader: &cachecontrolheader.Header{},
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := cachecontrolheader.Parse(tt.header)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.wantHeader, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
