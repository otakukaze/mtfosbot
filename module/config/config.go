package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"git.trj.tw/golang/mtfosbot/module/utils"
	yaml "gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	Port      string `yaml:"port"`
	URL       string `yaml:"url"`
	ImageRoot string `yaml:"image_root"`
	Line      struct {
		secret string `yaml:"secret"`
		access string `yaml:"access"`
	} `yaml:"line"`
}

var conf *Config

// LoadConfig -
func LoadConfig(p ...string) error {
	var fp string
	if len(p) > 0 && len(p[0]) > 0 {
		fp = p[0]
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		fp = path.Join(wd, "config.yml")
	}
	fp = utils.ParsePath(fp)

	exists := utils.CheckExists(fp, false)
	if !exists {
		return errors.New("config file not exists")
	}

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}

	conf = &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}

	return nil
}
