package config

import "github.com/mitchellh/mapstructure"

var DefaultFileName = "breeze.yaml"

type BuildConfiguration map[string]interface{}

func (c BuildConfiguration) GetPreset() string {
	if preset, ok := c["preset"]; ok {
		return preset.(string)
	}

	return ""
}

func (c BuildConfiguration) ToLaravel() (*LaravelBuildConfiguration, error) {
	config := &LaravelBuildConfiguration{}
	err := mapstructure.Decode(c, config)

	return config, err
}

type LaravelBuildConfiguration struct {
	Api bool `yaml:"api"`
}

type EnvironmentConfiguration struct {
	Image  string
	Domain string
}

type DeployConfiguration struct {
	Path string
}

type BreezeConfiguration struct {
	Project      string
	Build        BuildConfiguration
	Deploy       *DeployConfiguration
	Environments map[string]*EnvironmentConfiguration
}
