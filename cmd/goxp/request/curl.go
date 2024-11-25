package request

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/whitekid/goxp/requests"
)

func Run(ctx context.Context, url string) error {
	resp, err := requests.Get(url).
		Header(requests.HeaderUserAgent, *flag.userAgent).
		Do(ctx)
	if err != nil {
		return err
	}

	if *flag.verbose {
		fmt.Printf("%s\n", resp.Status)
		for k := range resp.Header {
			fmt.Printf("%s: %s\n", k, resp.Header.Get(k))
		}
		fmt.Printf("\n")
	}

	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return err
}
