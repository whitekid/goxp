package flags

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Int(fs *pflag.FlagSet, key string, name string, shorthand string, value int, usage string) (r *int) {
	r = fs.IntP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return
}

func Bool(fs *pflag.FlagSet, key string, name string, shorthand string, value bool, usage string) (r *bool) {
	r = fs.BoolP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return
}

func String(fs *pflag.FlagSet, key string, name string, shorthand string, value string, usage string) (r *string) {
	r = fs.StringP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return
}

func StringArray(fs *pflag.FlagSet, key string, name string, shorthand string, value []string, usage string) *[]string {
	r := fs.StringArrayP(name, shorthand, nil, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return r
}

func Duration(fs *pflag.FlagSet, key string, name string, shorthand string, value time.Duration, usage string) *time.Duration {
	r := fs.DurationP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return r
}
