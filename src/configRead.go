package src

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"sort"
)

type Config struct {
	City struct{
		Path string `yaml:"path"`
		SupportedLanguages []string `yaml:"supportedLanguages"`
		DefaultLanguage string `yaml:"defaultLanguage"`
	} `yaml:"city"`
	Asn struct{
		Path string `yaml:"path"`
	} `yaml:"asn"`
}

func (cfg *Config) CheckPath() error  {
	cleanPath := filepath.Clean(cfg.City.Path)
	if cleanPath != cfg.City.Path {
		return fmt.Errorf("the path of the City.mmdb file not standardized")
	}
	if cfg.City.Path == "" {
		return fmt.Errorf("the path of the City.mmdb file is empty")
	}

	cleanPath = filepath.Clean(cfg.Asn.Path)
	if cleanPath != cfg.Asn.Path {
		return fmt.Errorf("the path of the ASN.mmdb file not standardized")
	}
	if cfg.Asn.Path == "" {
		return fmt.Errorf("the path of the ASN.mmdb file is empty")
	}

	return nil
}

func (cfg *Config) CheckDefaultLanguage() error {
	sort.Strings(cfg.City.SupportedLanguages)
	index := sort.SearchStrings(cfg.City.SupportedLanguages, cfg.City.DefaultLanguage)
	if index > len(cfg.City.SupportedLanguages) || cfg.City.SupportedLanguages[index] != cfg.City.DefaultLanguage {
		return fmt.Errorf("the default language is not among the supported languages")
	}
	return nil
}

func (cfg *Config) Check() error {
	err := cfg.CheckPath()
	if err != nil {
		return err
	}
	err = cfg.CheckDefaultLanguage()
	if err != nil {
		return err
	}
	return  nil
}

func GetConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	cfg := &Config{}
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, cfg.Check()
}