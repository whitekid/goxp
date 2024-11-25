package request

import (
	"github.com/spf13/pflag"

	"github.com/whitekid/goxp/flags"
)

var (
	flag = struct {
		userAgent *string
		verbose   *bool
	}{}
)

func SetFlags(pfs *pflag.FlagSet, fs *pflag.FlagSet, version string) {
	flag.userAgent = flags.String(fs, "request.user-agent", "user-agent", "A", "goxp requests "+version, "use agent")
	flag.verbose = flags.Bool(fs, "request.verbose", "verbose", "v", false, "verbose")
}
