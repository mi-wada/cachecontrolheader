package cachecontrolheader

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		header     string
		wantHeader *Header
	}{
		{
			header: "max-age=3600, private, must-revalidate",
			wantHeader: &Header{
				MaxAge:         3600 * time.Second,
				Private:        true,
				MustRevalidate: true,
			},
		},
		{
			header:     "",
			wantHeader: &Header{},
		},
	} {
		t.Run(tt.header, func(t *testing.T) {
			h, err := Parse(tt.header)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.wantHeader, h); diff != "" {
				t.Errorf("Header mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
