package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMimeByExt(t *testing.T) {
	tests := [...]struct {
		name string
		want string
	}{
		{".json", "application/json; charset=utf-8"},
		{".png", "image/png"},
		{".gif", "image/gif"},
		{".vcf", "text/vcard; charset=utf-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, mimeByExt(tt.name))
		})
	}
}
