package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
)

var GlobalConfig = &GlobalConfiguration{}

type GlobalConfiguration struct {
	AuthToken string `json:"authToken"`
}

func (c *GlobalConfiguration) Load() error {
	f, err := os.Open(rootBreezeConfigPath())
	defer f.Close()

	if err == nil {
		// load the file
		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		yaml.Unmarshal(b, GlobalConfig)
		return nil
	}

	return err
}

func (c *GlobalConfiguration) Persist() {
	os.MkdirAll(rootBreezeDir(), 0777)

	f, _ := os.Create(rootBreezeConfigPath())
	defer f.Close()

	b, _ := yaml.Marshal(GlobalConfig)

	f.Write(b)
}

func rootBreezeDir() string {
	homedir := os.Getenv("HOME")

	return fmt.Sprintf("%s/.breeze", homedir)
}

func rootBreezeConfigPath() string {
	return fmt.Sprintf("%s/config.yaml", rootBreezeDir())
}
