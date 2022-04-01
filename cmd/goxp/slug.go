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
			fx.ForEach(
				fx.Times(10, func(i int) uuid.UUID { return uuid.New() }),
				func(i int, x uuid.UUID) {
					sg := slug.NewUUID()
					slug := sg.Encode(x)
					uid1 := sg.Decode(slug)

					fmt.Printf("%s => %s => %s\n", x, slug, uid1)
				})
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

			fx.ForEach(
				fx.Times(10, func(i int) *big.Int { return big.NewInt(int64(i)) }),
				func(i int, b *big.Int) {
					n := max.Int64() + int64(i)
					code := shortner.Encode(n)
					decoded, _ := shortner.Decode(code)
					fmt.Printf("%d => %s => %d\n", n, code, decoded)
				})
		},
	})

	rootCmd.AddCommand(cmd)
}
