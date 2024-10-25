package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type aesCipher struct {
	key []byte
}

var _ Interface = (*aesCipher)(nil)

func NewAes(key []byte) Interface { return &aesCipher{key: key} }

func (c *aesCipher) newCipher(key []byte) (cipher.Block, error) {
	if len(key) < aes.BlockSize {
		key = append(key, make([]byte, aes.BlockSize-len(key))...)
	}

	return aes.NewCipher(key)
}

func (c *aesCipher) Encrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("encrypt failed: %w", err)
	}

	if mod := len(data) % block.BlockSize(); mod != 0 {
		padding := make([]byte, block.BlockSize()-mod)
		data = append(data, padding...)
	}

	ciphertext := make([]byte, block.BlockSize()+len(data))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("random read failed: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], data)

	return ciphertext, nil
}

func (c *aesCipher) Decrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, errors.Errorf("data size mismatch")
	}

	iv := data[:block.BlockSize()]
	data = data[block.BlockSize():]

	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	plaintext = bytes.TrimRightFunc(plaintext, func(r rune) bool { return r == 0 })

	return plaintext, nil
}
