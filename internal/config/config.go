package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TeleConfig struct {
	Name   string
	Token  string
	ChatId int64
}

type SenderConfig struct {
	Name       string
	Tele       string
	Folder     string
	TeleConfig *TeleConfig
}

type MasterConfig struct {
	Senders     []SenderConfig
	TeleConfigs []TeleConfig
}

func validateConfig(cfg MasterConfig) error {
	return nil
}

func ReadConfig(path_to_config string) (MasterConfig, error) {

	config := MasterConfig{}
	var yfile []byte
	{
		var err error
		yfile, err = os.ReadFile(path_to_config)
		if err != nil {
			return config, err
		}
	}

	{
		err := yaml.Unmarshal(yfile, &config)
		if err != nil {
			return config, err
		}
	}

	{
		for si, sender := range config.Senders {
			for ti, tele := range config.TeleConfigs {
				if tele.Name == sender.Tele {
					config.Senders[si].TeleConfig = &config.TeleConfigs[ti]
				}
			}
		}
	}

	{
		err := validateConfig(config)
		if err != nil {
			return config, err
		}
	}

	return config, nil
}
