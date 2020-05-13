package config

var DefaultFileName = "breeze.yaml"

type BuildConfiguration struct {
	Preset string
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
	Build        *BuildConfiguration
	Deploy       *DeployConfiguration
	Environments map[string]*EnvironmentConfiguration
}
