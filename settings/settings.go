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
	Rotations
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

func (s Setup) Submittable() bool {
	strs := s.Name != "" &&
		s.Source != ""
	ints := s.Frequency >= 1 &&
		s.Frequency <= 1440 &&
		s.Initial >= 1 &&
		(s.Daily <= 28 ||
			s.Daily > 300)
	dest := func() (d bool) {
		if s.Type == "remote" {
			d = (s.LocalHost != "" ||
				s.RemoteHost != "") &&
				s.Username != "" &&
				s.Path != "" &&
				s.Port != "" &&
				s.PrivateKeyUrl != ""
		}
		if s.Type == "mount" {
			d = s.Path != ""
		}
		return d
	}
	return strs && ints && dest()
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
	newConf := filepath.Join(configDir, answerSet.Name+".toml")
	f, err := os.Create(newConf)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var config = map[string]interface{}{
		"name":        answerSet.Name,
		"description": answerSet.Description,
		"source":      answerSet.Source,
		"excludes":    answerSet.Excludes,
		"destination": answerSet.Destination,
		"rotations":   answerSet.Rotations,
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
