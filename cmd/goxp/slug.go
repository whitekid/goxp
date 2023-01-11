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
			fx.Each(
				fx.Times(10, func(i int) uuid.UUID { return uuid.New() }),
				func(i int, x uuid.UUID) {
					sg := slug.NewUUID(nil)
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
			encoding := string(fx.Shuffle([]byte(slug.EncodeURL)))
			shortner := slug.NewShortner(encoding)
			fmt.Printf("encoding: %s\n\n", encoding)

			max, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt16))

			fx.Each(
				fx.Times(10, func(i int) *big.Int { return big.NewInt(int64(i)) }),
				func(i int, b *big.Int) {
					n := max.Int64() + int64(i)
					code := shortner.Encode(n)
					decoded, _ := shortner.Decode(code)
					fmt.Printf("%d => %s => %d\n", n, code, decoded)
				})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "new",
		Short: "generate new random encoding",
		Run: func(cmd *cobra.Command, args []string) {
			enc := fx.Shuffle([]rune(slug.EncodeURL))
			fmt.Printf("%s\n", string(enc))
		},
	})

	rootCmd.AddCommand(cmd)
}
