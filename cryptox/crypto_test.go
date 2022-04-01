package cryptox

import (
	"crypto/aes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
)

func TestCrypt(t *testing.T) {
	type args struct {
		key  string
		data string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"default", args{
			key:  goxp.RandomString(aes.BlockSize),
			data: `동해 물과 백두산이 마르고 닳도록 동해 물과 백두산이 마르고 닳도록`,
		}},
		{"random", args{
			key:  goxp.RandomString(aes.BlockSize),
			data: goxp.RandomString(rand.Intn(1024)),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.args.key, tt.args.data)
			require.NoError(t, err)
			require.NotEqual(t, tt.args.data, encrypted)

			decrypted, err := Decrypt(tt.args.key, encrypted)
			require.NoError(t, err)

			require.Equal(t, tt.args.data, string(decrypted))

			{
				c := NewAes([]byte(tt.args.key))
				encrypted, err := c.Encrypt([]byte(tt.args.data))
				require.NoError(t, err)
				require.NotEqual(t, tt.args.data, encrypted)

				decrypted, err := c.Decrypt(encrypted)
				require.NoError(t, err)
				require.Equal(t, tt.args.data, string(decrypted))
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	key := goxp.RandomString(aes.BlockSize)
	data := "hello_world"

	enc1 := MustEncrypt(key, data)
	enc2 := MustEncrypt(key, data)
	require.NotEqual(t, enc1, enc2)
}

func TestEncryptError(t *testing.T) {
	type args struct {
		key  []byte
		data []byte
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"small key", args{
			key:  []byte("key"),
			data: []byte("value"),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAes([]byte(tt.args.key)).Encrypt(tt.args.data)
			require.NoError(t, err)
		})
	}
}
