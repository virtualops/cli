package config

var DefaultFileName = "breeze.yaml"

type BuildConfiguration struct {
	Preset string
}

type EnvironmentConfiguration struct {
	Path string
}

type BreezeConfiguration struct {
	Project      string
	Build        *BuildConfiguration
	Environments map[string]*EnvironmentConfiguration
}
