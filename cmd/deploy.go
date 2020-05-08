package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/virtualops/breeze-cli/pkg/http"
	"github.com/virtualops/breeze-cli/pkg/projects"
	"os"
	"time"
)

var deployCmd = &cobra.Command{
	Use: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := projects.GetProject(Config.Project)

		if err != nil && http.Is404(err) {
			fmt.Sprintf("The project %s doesn't exist\n", Config.Project)
			os.Exit(1)
		}

		var environmentName string

		for _, env := range project.Environments {
			if env.Type != "production" {
				environmentName = env.Name
			}
		}

		if len(args) > 0 {
			environmentName = args[0]
		} else if environmentName == "" {
			fmt.Println("No test environments found. If you would like to deploy to production, try breeze deploy production")
		}

		s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = " Building docker image"
		s.Start()
		time.Sleep(14 * time.Second)
		s.Stop()

		fmt.Println("\033[1;32m✅ Built image\033[0m")

		s = spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = " Uploading image to AWS"
		s.Start()
		time.Sleep(4 * time.Second)
		s.Stop()

		fmt.Println("\033[1;32m✅ Uploaded to AWS\033[0m")

		s = spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = " Building Upgrade Diff"
		s.Start()
		time.Sleep(1 * time.Second)
		s.Stop()

		s = spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Deploying to %s", environmentName)
		s.Start()
		time.Sleep(12 * time.Second)
		s.Stop()

		fmt.Println("\033[1;32m✅ Deploy completed\033[0m")
	},
}
