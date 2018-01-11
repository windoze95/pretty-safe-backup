package run

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	RunConfigDir string `toml:"run-config-directory"`
	logDir       string
}

func (c *Config) initDefaults() {
	c.RunConfigDir = "/etc/xdg/psb/run"
	c.logDir = "/var/log/psb/"
}

var config = Config{}

func init() {
	config.initDefaults()
	configFile := "/etc/xdg/psb/config.toml"
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			log.Fatal(err)
		}
	}
	if err := os.MkdirAll(config.RunConfigDir, 0755); err != nil {
		log.Fatalf("error creating directory %s: %s", config.RunConfigDir, err.Error())
	}
}
