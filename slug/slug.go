package slug

import (
	"encoding/base64"

	"github.com/google/uuid"
)

const EncodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_" // url encoding; base64.encodeURL 와 동일

type Slug struct {
	encoder *base64.Encoding
}

func New(encoder string) *Slug {
	return withEncoding(base64.NewEncoding(encoder).WithPadding(base64.NoPadding))
}

func withEncoding(encoding *base64.Encoding) *Slug {
	return &Slug{
		encoder: encoding,
	}
}

func (slug *Slug) Encode(src []byte) string        { return slug.encoder.EncodeToString(src) }
func (slug *Slug) Decode(s string) ([]byte, error) { return slug.encoder.DecodeString(s) }

type UUID struct {
	slugger *Slug
}

func NewUUID(encoding *base64.Encoding) *UUID {
	if encoding == nil {
		encoding = base64.RawURLEncoding
	}

	return &UUID{
		slugger: withEncoding(encoding),
	}
}

// Encode uuid to slug
func (s *UUID) Encode(uid uuid.UUID) string {
	return s.slugger.Encode(uid[:])
}

// Decode slug to uuid
func (s *UUID) Decode(slug string) uuid.UUID {
	dec, err := s.slugger.Decode(slug)
	if err != nil {
		return uuid.UUID{}
	}

	if uid, err := uuid.FromBytes(dec); err == nil {
		return uid
	}

	return uuid.UUID{}
}

var uuidSlugger = NewUUID(nil)

// ToSlug ...
// Deprecated:
func ToSlug(uid uuid.UUID) string { return uuidSlugger.Encode(uid) }

// ToUUID ...
// Deprecated:
func ToUUID(slug string) uuid.UUID { return uuidSlugger.Decode(slug) }
