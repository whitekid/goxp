package slug

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {
	uid := uuid.New()

	slug := Slug(uid)
	uid1 := UUID(slug)

	require.Equal(t, uid.String(), uid1.String())

	log.Printf("uuid: %s, slug=%s", uid.String(), slug)
}

func TestExample(t *testing.T) {
	uid := uuid.New()

	slug := Slug(uid)
	fmt.Printf("uuid=%s\n", uid.String())
	fmt.Printf("slug=%s\n", slug)

	uid1 := UUID(slug)
	fmt.Printf("decode=%s\n", uid1.String())
}
