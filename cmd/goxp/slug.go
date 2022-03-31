package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/slug"
)

func init() {
	cmd := &cobra.Command{
		Use:   "slug",
		Short: "slug package examples",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "uuid",
		Short: "encode uuid to URL friendly",
		Run: func(cmd *cobra.Command, args []string) {
			uid := uuid.New()
			sg := slug.NewUUID()
			slug := sg.Encode(uid)

			fmt.Printf("UUID: %s => slug=%s\n", uid, slug)

			uid1 := sg.Decode(slug)
			fmt.Printf("slug: %s => UUID=%s\n", slug, uid1)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "short",
		Short: "encode int to URL friendly.",
		Run: func(cmd *cobra.Command, args []string) {
			encoding := fx.Shuffle([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"))
			shortner := slug.NewShortner(string(encoding))
			fmt.Printf("encoding: %s\n\n", string(encoding))

			max, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt16))

			for i := max.Int64(); i < max.Int64()+10; i++ {
				fmt.Printf("%d => %s\n", i, shortner.Encode(i))
			}
		},
	})

	rootCmd.AddCommand(cmd)
}
