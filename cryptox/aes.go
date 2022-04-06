package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pkg/errors"
)

type aesCipher struct {
	key []byte
}

func NewAes(key []byte) Interface {
	return &aesCipher{key: key}
}

func (c *aesCipher) newCipher(key []byte) (cipher.Block, error) {
	if len(key) < aes.BlockSize {
		key = append(key, make([]byte, aes.BlockSize-len(key))...)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "newCipher failed")
	}

	return block, nil
}

func (c *aesCipher) Encrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, errors.Wrap(err, "encrypt failed")
	}

	if mod := len(data) % block.BlockSize(); mod != 0 {
		padding := make([]byte, block.BlockSize()-mod)
		data = append(data, padding...)
	}

	ciphertext := make([]byte, block.BlockSize()+len(data))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Wrapf(err, "random read failed")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], data)

	return ciphertext, nil
}

func (c *aesCipher) Decrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, err
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, errors.Errorf("data size mismatch")
	}

	iv := data[:block.BlockSize()]
	data = data[block.BlockSize():]

	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	plaintext = bytes.TrimRightFunc(plaintext, func(r rune) bool { return byte(r) == byte(00) })

	return plaintext, nil
}

var _ Interface = (*aesCipher)(nil)
