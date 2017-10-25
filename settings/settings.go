package settings

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/casimir/xdg-go"
)

type Setup struct {
	Name        string
	Description string
	Source      string
	Excludes    []string
	Destination
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

func (s Setup) Submittable() bool {
	strs := s.Name != "" &&
		s.Source != ""
	maps := func() bool {
		dest := ((s.LocalHost != "" ||
			s.RemoteHost != "") &&
			s.Username != "" &&
			s.Path != "" &&
			s.Port != "" &&
			s.PrivateKeyUrl != "") ||
			s.Path != ""
		return dest
	}
	return strs && maps()
}

var (
	App       = xdg.App{Name: "psb"}
	configDir = filepath.Join(xdg.ConfigHome(), App.Name)
)

func init() {
	err := os.MkdirAll(configDir, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func NewConfig(answerSet *Setup) {
	newOperationPath := filepath.Join(configDir, answerSet.Name)
	err := os.MkdirAll(newOperationPath, 0777)
	if err != nil {
		log.Fatal(err)
	}
	newConf := filepath.Join(newOperationPath, "rc.toml")
	f, err := os.Create(newConf)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var config = map[string]interface{}{
		"description": answerSet.Description,
		"source":      answerSet.Source,
		"excludes":    answerSet.Excludes,
		"destination": answerSet.Destination,
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(f, strings.NewReader(buf.String()))
	if err != nil {
		log.Fatal(err)
	}
}
