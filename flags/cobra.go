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
	Key          string
	Short        string
	DefaultValue interface{}
	Description  string
}

var configs map[string][]Flag

// InitDefaults initialize config
func InitDefaults(configs map[string][]Flag) {
	for use := range configs {
		for _, config := range configs[use] {
			if config.DefaultValue != nil {
				viper.SetDefault(config.Key, config.DefaultValue)
			}
		}
	}
}

// InitFlagSet init flags to cobra command
func InitFlagSet(configs map[string][]Flag, use string, fs *pflag.FlagSet) {
	uses := strings.Split(use, " ")
	for _, config := range configs[uses[0]] {
		switch v := config.DefaultValue.(type) {
		case int:
			fs.IntP(config.Key, config.Short, v, config.Description)
		case bool:
			fs.BoolP(config.Key, config.Short, v, config.Description)
		case string:
			fs.StringP(config.Key, config.Short, v, config.Description)
		case []byte:
			fs.BytesHexP(config.Key, config.Short, v, config.Description)
		case time.Duration:
			fs.DurationP(config.Key, config.Short, v, config.Description)
		default:
			log.Errorf("unsupported type %T", config.DefaultValue)
		}
		viper.BindPFlag(config.Key, fs.Lookup(config.Key))
	}
}
