package settings

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/orange-lightsaber/pretty-safe-backup/util"
)

const (
	runConfigDir       = "/etc/xdg/psb/run"
	PathFromTimeLayout = "/2006/January/2/1504Z"
)

type RunConfig struct {
	enabled     bool
	Name        string
	Description string
	Source      string
	Excludes    []string
	Destination Destination
	Rotations   Rotations
}

type Destination struct {
	Path          string
	LocalHost     string
	RemoteHost    string
	Username      string
	Port          string
	PrivateKeyUrl string
	Type          string
}

type Rotations struct {
	Frequency int
	Initial   int
	Daily     int
	Monthly   int
	Yearly    int
}

func (c *RunConfig) EnableConfig() {
	c.enabled = true
}

func (c *RunConfig) DisableConfig() {
	c.enabled = false
}

func (c RunConfig) Submittable() bool {
	head := c.Name != "" &&
		c.Source != ""
	rotas := c.Rotations.Frequency >= 1 &&
		c.Rotations.Frequency <= 1440 &&
		c.Rotations.Initial >= 1 &&
		(c.Rotations.Daily <= 28 ||
			c.Rotations.Daily > 300)
	dest := func() (d bool) {
		if c.Destination.Type == "remote" {
			d = (c.Destination.LocalHost != "" ||
				c.Destination.RemoteHost != "") &&
				c.Destination.Username != "" &&
				c.Destination.Path != "" &&
				c.Destination.Port != "" &&
				c.Destination.PrivateKeyUrl != ""
		}
		if c.Destination.Type == "mount" {
			d = c.Destination.Path != ""
		}
		return d
	}
	return head && rotas && dest()
}

func init() {
	if err := os.MkdirAll(runConfigDir, 0755); err != nil {
		log.Fatal(util.ErrorMsg("Did you forget to use sudo? :)", err))
	}
}

func WriteRunConfig(runConfig *RunConfig) {
	newRunConfigFile := filepath.Join(runConfigDir, runConfig.Name+".toml")
	f, err := os.Create(newRunConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	c := map[string]interface{}{
		"enabled":     runConfig.enabled,
		"name":        runConfig.Name,
		"description": runConfig.Description,
		"source":      runConfig.Source,
		"excludes":    runConfig.Excludes,
		"destination": runConfig.Destination,
		"rotations":   runConfig.Rotations,
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(f, strings.NewReader(buf.String())); err != nil {
		log.Fatal(err)
	}
}

func EnableRunConfig(name string) error {
	errChan := make(chan error)
	runConfigChan := make(chan RunConfig)
	go GetDisabledConfigs(runConfigChan, errChan)
	for c := range runConfigChan {
		if c.Name == name {
			c.EnableConfig()
			WriteRunConfig(&c)
		}
	}
	return <-errChan
}

func DisableRunConfig(name string) error {
	errChan := make(chan error)
	runConfigChan := make(chan RunConfig)
	go GetEnabledConfigs(runConfigChan, errChan)
	for c := range runConfigChan {
		if c.Name == name {
			c.DisableConfig()
			WriteRunConfig(&c)
		}
	}
	return <-errChan
}

func GetEnabledConfigs(runConfigChan chan RunConfig, errChan chan error) {
	getConfigs(runConfigChan, true, errChan)
}

func GetDisabledConfigs(runConfigChan chan RunConfig, errChan chan error) {
	getConfigs(runConfigChan, false, errChan)
}

func getConfigs(runConfigChan chan RunConfig, enabled bool, errChan chan error) {
	var wg sync.WaitGroup
	err := filepath.Walk(runConfigDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if c, err := decodeRunConfig(path); err != nil {
				errChan <- err
			} else {
				runConfigChan <- c
			}
		}()
		return err
	})
	wg.Wait()
	if err != nil {
		errChan <- util.ErrorMsg("error parsing run config directory (was psb ran by a superuser?)", err)
	}
	close(runConfigChan)
}

// TODO create Config interface to impliment global config
func decodeRunConfig(path string) (RunConfig, error) {
	config := struct {
		Enabled bool
		RunConfig
	}{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		err = util.ErrorMsg("error decoding run config", err)
	}
	c := config.RunConfig
	// c.Format()
	if config.Enabled {
		c.EnableConfig()
	}
	return c, err
}
