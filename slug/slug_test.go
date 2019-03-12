package slug

import (
	"fmt"
	"log"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestSlug(t *testing.T) {
	uid := uuid.NewV4()

	slug := Slug(uid)
	uid1 := UUID(slug)

	assert.Equal(t, uid.String(), uid1.String())

	log.Printf("uuid: %s, slug=%s", uid.String(), slug)
}

func TestExample(t *testing.T) {
	uid := uuid.NewV4()

	slug := Slug(uid)
	fmt.Printf("uuid=%s\n", uid.String())
	fmt.Printf("slug=%s\n", slug)

	uid1 := UUID(slug)
	fmt.Printf("decode=%s\n", uid1.String())
}
