package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/virtualops/cli/pkg/config"
	"io/ioutil"
	"os"
)

func loadGlobalConfig(command *cobra.Command, args []string) {
	homedir := os.Getenv("HOME")
	configFilePath := fmt.Sprintf("%s/.breeze/config.yaml", homedir)
	f, err := os.Open(configFilePath)

	if err == nil {
		// load the file
		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		yaml.Unmarshal(b, config.GlobalConfig)
		return
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	os.MkdirAll(fmt.Sprintf("%s/.breeze", homedir), 0777)

	f, _ = os.Create(configFilePath)
	config.GlobalConfig = &config.GlobalConfiguration{}
	b, _ := yaml.Marshal(config.GlobalConfig)

	f.Write(b)
}
