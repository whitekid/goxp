package slug

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := [...]struct {
		name     string
		args     args
		wantSlug string
	}{
		{"default", args{"fcf9853b-27aa-4ea1-b60c-2b2f443afb1a"}, "_PmFOyeqTqG2DCsvRDr7Gg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := uuid.Parse(tt.args.uuid)
			require.NoError(t, err)

			slug := ToSlug(uid)
			require.Equal(t, tt.wantSlug, slug)

			uid1 := ToUUID(slug)
			require.Equal(t, uid.String(), uid1.String())
		})
	}
}

func TestSlugger(t *testing.T) {
	sl := New(EncodeURL)

	got := sl.Encode(big.NewInt(1).Bytes())
	require.Equal(t, "AQ", got)

	dec, err := sl.Decode(got)
	require.NoError(t, err)

	b := big.NewInt(1)

	require.Equal(t, b.Bytes(), dec)

	s := New(EncodeURL)
	require.Equal(t, "AQ", s.Encode(big.NewInt(1).Bytes()))
}
