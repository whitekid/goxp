package flags

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Int32(flags *pflag.FlagSet, key string, name string, shorthand string, value int32, usage string) (r *int32) {
	r = flags.Int32P(name, shorthand, value, usage)
	viper.BindPFlag(key, flags.Lookup(name))

	return
}

func Bool(flags *pflag.FlagSet, key string, name string, shorthand string, value bool, usage string) (r *bool) {
	r = flags.BoolP(name, shorthand, value, usage)
	viper.BindPFlag(key, flags.Lookup(name))

	return
}

func String(flags *pflag.FlagSet, key string, name string, shorthand string, value string, usage string) (r *string) {
	r = flags.StringP(name, shorthand, value, usage)
	viper.BindPFlag(key, flags.Lookup(name))

	return
}

func StringArray(flags *pflag.FlagSet, key string, name string, shorthand string, value []string, usage string) *[]string {
	r := flags.StringArrayP(name, shorthand, nil, usage)
	viper.BindPFlag(key, flags.Lookup(name))
	return r
}
