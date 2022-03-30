package slug

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func New() string {
	return ToSlug(uuid.New())
}

// ToSlug uuid to slug
func ToSlug(uid uuid.UUID) string {
	encoded := base64.URLEncoding.EncodeToString(uid[0:])
	return encoded[:22]
}

// ToUUID slug to uuid
func ToUUID(slug string) uuid.UUID {
	data, err := base64.URLEncoding.DecodeString(slug + "==")
	if err != nil {
		return uuid.UUID{}
	}

	uid, err := uuid.FromBytes(data)
	if err != nil {
		return uuid.UUID{}
	}
	return uid
}
