package config

import (
	"encoding/json"
	"io/ioutil"
)

var AppConfig *Config

type Config struct {
	DSN   string `json:"DSN"`
	Table string `json:"table"`

	HttpPort         int `json:"httpPort"`
	HttpReadTimeout  int `json:"httpReadTimeout"`
	HttpWriteTimeout int `json:"httpWriteTimeout"`
}

func LoadConf(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	conf := Config{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return err
	}
	AppConfig = &conf
	return nil
}
