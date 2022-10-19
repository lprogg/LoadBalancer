package util

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func LoadConfig(r io.Reader) (*Config, error) {
	buf, err := ioutil.ReadAll(r)

	if err != nil {
		return nil, err
	}

	conf := Config{}

	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
