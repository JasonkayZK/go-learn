package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	appConfig map[string]string
)

func init() {
	err := loadConfig()
	if err != nil {
		fmt.Printf("failed to load config")
	}
}

func loadConfig() error {
	appConfig = map[string]string{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to load config")
	}

	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config")
	}

	return nil
}

func GetConfig(key string) string {
	if value, exist := appConfig[key]; exist {
		return value
	} else {
		return "error"
	}
}
