package request

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/whitekid/goxp/flags"
)

const (
	fkUserAgent = "request.user-agent"
	fkVerbose   = "request.verbose"
)

func SetFlags(pfs *pflag.FlagSet, fs *pflag.FlagSet, version string) {
	flags.String(fs, fkUserAgent, "user-agent", "A", "goxp requests "+version, "use agent")
	flags.Bool(fs, fkVerbose, "verbose", "v", false, "verbose")
}

func userAgent() string { return viper.GetString(fkUserAgent) }
func verbose() bool     { return viper.GetBool(fkVerbose) }
