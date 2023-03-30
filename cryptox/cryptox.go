package cryptox

import (
	"crypto/aes"
	"encoding/base64"

	"github.com/whitekid/goxp/log"
)

type Encrypter interface {
	Encrypt(data []byte) ([]byte, error)
}

type Decrypter interface {
	Decrypt(data []byte) ([]byte, error)
}

type Interface interface {
	Encrypter
	Decrypter
}

const BlockSize = aes.BlockSize

func Encrypt(key, data string) (string, error) {
	encrypted, err := NewAes([]byte(key)).Encrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(key, data string) (string, error) {
	encrypted, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	decrypted, err := NewAes([]byte(key)).Decrypt(encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func MustEncrypt(key, data string) string {
	enc, err := Encrypt(key, data)
	if err != nil {
		log.Fatalf("%+v key=%s data=%s", err, key, data)
	}
	return enc
}

func MustDecrypt(key, data string) string {
	dec, err := Decrypt(key, data)
	if err != nil {
		log.Fatalf("%+v key=%s data=%s", err, key, data)
	}
	return dec
}
