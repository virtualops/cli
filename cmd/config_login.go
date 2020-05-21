package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/virtualops/cli/pkg/api"
	"github.com/virtualops/cli/pkg/config"
	"os"
	"regexp"
)

var emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var forceLogin bool

var loginCmd = &cobra.Command{
	Use: "login",
	PersistentPreRun: loadGlobalConfig,
	Run: func(cmd *cobra.Command, args []string) {
		if config.GlobalConfig.AuthToken != "" && !forceLogin {
			fmt.Printf("%s You are already logged in\n", promptui.IconGood)
			os.Exit(0)
		}

		for {
			prompt := promptui.Prompt{
				Label: "Email",
				Validate: func(s string) error {
					re := regexp.MustCompile(emailRegex)
					if !re.MatchString(s) {
						return errors.New("that email address is not valid")
					}

					return nil
				},
			}

			email, _ := prompt.Run()
			prompt = promptui.Prompt{
				Label: "Password",
				Mask: '•',
			}
			password, _ := prompt.Run()

			token, err := api.Api.Login(email, password)

			if err != nil {
				continue
			}

			fmt.Printf("%s Logged in as %s\n", promptui.IconGood, email)
			config.GlobalConfig.AuthToken = token
			break
		}
		config.GlobalConfig.Persist()
	},
}

func init() {
	loginCmd.Flags().BoolVarP(&forceLogin, "force", "f", false, "Forcibly re-login, even if already authenticated")
}
