// Package flags provides flags defaults frameworks
package flags

import (
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/whitekid/go-utils/log"
)

type Flag struct {
	Name         string
	Shorthand    string
	DefaultValue interface{}
	Usage        string
}

var configs map[string][]Flag

// InitDefaults initialize config
func InitDefaults(v *viper.Viper, configs map[string][]Flag) {
	if v == nil {
		v = viper.GetViper()
	}

	for use := range configs {
		for _, config := range configs[use] {
			if config.DefaultValue != nil {
				v.SetDefault(config.Name, config.DefaultValue)
			}
		}
	}
}

// InitFlagSet init flags to cobra command
func InitFlagSet(v *viper.Viper, configs map[string][]Flag, use string, fs *pflag.FlagSet) {
	if v == nil {
		v = viper.GetViper()
	}

	uses := strings.Split(use, " ")
	for _, cfg := range configs[uses[0]] {
		switch v := cfg.DefaultValue.(type) {
		case int:
			fs.IntP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		case bool:
			fs.BoolP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		case string:
			fs.StringP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		case []string:
			fs.StringSliceP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		case []byte:
			fs.BytesHexP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		case time.Duration:
			fs.DurationP(cfg.Name, cfg.Shorthand, v, cfg.Usage)
		default:
			log.Errorf("unsupported type %T", cfg.DefaultValue)
		}

		var flag *pflag.Flag
		if cfg.Name != "" {
			flag = fs.Lookup(cfg.Name)
		}

		if flag == nil && cfg.Shorthand != "" {
			flag = fs.ShorthandLookup(cfg.Shorthand)
		}

		if flag == nil {
			log.Debugf("flag not found %+v", cfg)
			continue
		}

		v.BindPFlag(cfg.Name, flag)
	}
}
