package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStruct(t *testing.T) {
	type args struct {
		s any
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`required`, args{struct {
			X int `validate:"required"`
		}{}}, true},
		{`required`, args{struct {
			X int `validate:"required"`
		}{X: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.args.s)
			require.Truef(t, (err != nil) == tt.wantErr, `Struct() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
		})
	}
}

func TestVar(t *testing.T) {
	type args struct {
		v    any
		spec string
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`required`, args{0, "required"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Var(tt.args.v, tt.args.spec)
			require.Truef(t, (err != nil) == tt.wantErr, `Var() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	{
		err := Var(0, "required")
		require.Error(t, err)
		require.Truef(t, IsValidationError(err), "err: %T %s", err, err)
	}

	{
		err := Var(1, "required")
		require.NoError(t, err)
		require.False(t, IsValidationError(err))
	}
}
