package cryptox

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"io"

	"github.com/whitekid/goxp/errors"
)

type desCipher struct {
	key []byte
}

var _ Interface = (*desCipher)(nil)

func NewDes(key []byte) Interface { return &desCipher{key: key} }

func (c *desCipher) newCipher(key []byte) (cipher.Block, error) {
	if len(key) < des.BlockSize {
		key = append(key, make([]byte, des.BlockSize-len(key))...)
	}

	return des.NewCipher(key)
}

func (c *desCipher) Encrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, errors.Errorf(err, "encrypt failed")
	}

	if mod := len(data) % block.BlockSize(); mod != 0 {
		padding := make([]byte, block.BlockSize()-mod)
		data = append(data, padding...)
	}

	ciphertext := make([]byte, block.BlockSize()+len(data))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Errorf(err, "random read failed")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], data)

	return ciphertext, nil
}

func (c *desCipher) Decrypt(data []byte) ([]byte, error) {
	block, err := c.newCipher(c.key)
	if err != nil {
		return nil, errors.Errorf(err, "decrypt failed")
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("data size mismatch")
	}

	iv := data[:block.BlockSize()]
	data = data[block.BlockSize():]

	plaintext := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, data)

	plaintext = bytes.TrimRightFunc(plaintext, func(r rune) bool { return r == 0 })

	return plaintext, nil
}
