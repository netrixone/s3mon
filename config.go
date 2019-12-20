package main

import (
	"github.com/jroimartin/gocui"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

const (
	ResultErr   = "ERR"
	OkBgColor = gocui.ColorGreen
	OkFgColor = gocui.ColorBlack
	WarnBgColor = gocui.ColorYellow
	WarnFgColor = gocui.ColorBlack
	ErrBgColor  = gocui.ColorRed
	ErrFgColor  = gocui.ColorBlack

	defaultThreshold = float32(0.3)
)

type Config struct {
	Items []*Item `yaml:"items"`
}

type Item struct {
	Label     string  `yaml:"label,omitempty"`
	Script    string  `yaml:"script"`
	Unit      string  `yaml:"unit,omitempty"`
	Threshold float32 `yaml:"threshold,omitempty"`
	//Color     ui.Color `yaml:"color,omitempty"`
}

func NewConfig(path string) *Config {
	config := readFile(path)
	config.setDefaults()
	return config
}

func readFile(location string) *Config {
	yamlFile, err := ioutil.ReadFile(location)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", location)
	}

	config := new(Config)
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	return config
}

func (config *Config) setDefaults() {
	for _, item := range config.Items {
		if item.Threshold == 0 {
			item.Threshold = defaultThreshold
		}
	}
}
