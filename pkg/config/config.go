package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"path/filepath"
)

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
	Api bool `json:"api"`
}

type EnvironmentConfiguration struct {
	Image  string
	Domain string
}

type DeployConfiguration struct {
	Path  string
	Paths []string
}

type BreezeConfiguration struct {
	Project      string
	Build        BuildConfiguration
	Deploy       *DeployConfiguration
	Environments map[string]*EnvironmentConfiguration
}

func (c *BreezeConfiguration) ReleaseName() string {
	cwd, _ := os.Getwd()
	dir := filepath.Base(cwd)
	return fmt.Sprintf("%s-%s", c.Project, dir)
}

func (c *BreezeConfiguration) ImageName() string {
	cwd, _ := os.Getwd()
	dir := filepath.Base(cwd)
	return fmt.Sprintf("%s/%s", c.Project, dir)
}
