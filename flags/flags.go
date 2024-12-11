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
func IntVar(fs *pflag.FlagSet, p *int, key string, name string, shorthand string, value int, usage string) {
	fs.IntVarP(p, name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))
}

func Bool(fs *pflag.FlagSet, key string, name string, shorthand string, value bool, usage string) (r *bool) {
	r = fs.BoolP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return
}

func BoolVar(fs *pflag.FlagSet, p *bool, key string, name string, shorthand string, value bool, usage string) {
	fs.BoolVarP(p, name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))
}

func String(fs *pflag.FlagSet, key string, name string, shorthand string, value string, usage string) (r *string) {
	r = fs.StringP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return
}

func StringVar(fs *pflag.FlagSet, p *string, key string, name string, shorthand string, value string, usage string) {
	fs.StringVarP(p, name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))
}

func StringArray(fs *pflag.FlagSet, key string, name string, shorthand string, value []string, usage string) *[]string {
	r := fs.StringArrayP(name, shorthand, nil, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return r
}

func StringArrayVar(fs *pflag.FlagSet, p *[]string, key string, name string, shorthand string, value []string, usage string) {
	fs.StringArrayVarP(p, name, shorthand, nil, usage)
	viper.BindPFlag(key, fs.Lookup(name))
}

func Duration(fs *pflag.FlagSet, key string, name string, shorthand string, value time.Duration, usage string) *time.Duration {
	r := fs.DurationP(name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))

	return r
}

func DurationVar(fs *pflag.FlagSet, p *time.Duration, key string, name string, shorthand string, value time.Duration, usage string) {
	fs.DurationVarP(p, name, shorthand, value, usage)
	viper.BindPFlag(key, fs.Lookup(name))
}
