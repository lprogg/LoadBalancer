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

	config := Config{}

	if err := yaml.Unmarshal(buf, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
