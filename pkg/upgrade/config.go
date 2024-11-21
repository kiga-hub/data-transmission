package upgrade

import "github.com/spf13/viper"

const (
	dir = "transmission.dir"
)

var defaultConfig = Config{
	Dir: "/data/transmission",
}

// Config -
type Config struct {
	Dir string `toml:"dir" json:"dir,omitempty"`
}

// SetDefaultConfig -
func SetDefaultConfig() {
	viper.SetDefault(dir, defaultConfig.Dir)
}

// GetConfig -
func GetConfig() *Config {
	return &Config{
		Dir: viper.GetString(dir),
	}
}
