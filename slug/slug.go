package slug

import (
	"encoding/base64"

	"github.com/satori/go.uuid"
)

func Slug(uid uuid.UUID) string {
	encoded := base64.URLEncoding.EncodeToString(uid.Bytes())
	return encoded[:22]
}

func UUID(slug string) uuid.UUID {
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
