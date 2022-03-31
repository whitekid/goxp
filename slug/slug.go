package slug

import (
	"encoding/base64"

	"github.com/google/uuid"
)

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

func NewUUID() *UUID {
	return &UUID{
		slugger: withEncoding(base64.RawURLEncoding),
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

var uuidSlugger = NewUUID()

// ToSlug ...
// Deprecated:
func ToSlug(uid uuid.UUID) string { return uuidSlugger.Encode(uid) }

// ToUUID ...
// Deprecated:
func ToUUID(slug string) uuid.UUID { return uuidSlugger.Decode(slug) }
