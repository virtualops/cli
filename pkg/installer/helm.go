package installer

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"os/exec"
	"time"
)

func InstallOrVerifyHelm() {
	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	// check for tilt
	err = exec.Command("which", "helm").Run()

	if err != nil {
		fmt.Println("No Helm installation found, installing now")
		err = exec.Command("which", "brew").Run()

		// brew exists
		if err == nil {
			s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
			s.Suffix = " Installing Helm via Homebrew"
			s.Start()
			cmd := exec.Command("brew", "install", "helm")
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			s.Stop()
			if err != nil {
				fmt.Println("Failed to install Helm")
				os.Exit(1)
			}
		} else {
			// brew not available, so we'll use the install script (https://helm.sh/docs/intro/install/#from-script)
			installScriptSource := "https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3"
			err := exec.Command("curl", "-fsSL", installScriptSource, "-o", "get_helm.sh").Run()
			if err != nil {
				fmt.Println("Failed to download Helm installer")
				os.Exit(1)
			}
			s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
			s.Suffix = fmt.Sprintf(" Installing Helm using official install script (%s)", installScriptSource)
			s.Start()
			cmd := exec.Command("bash", "get_helm.sh")
			cmd.Stdout = os.Stdout
			err = cmd.Run()
			s.Stop()
			os.Remove(fmt.Sprintf("%s/%s", cwd, "get_helm.sh"))
			if err != nil {
				fmt.Println("Failed to install Helm")
				os.Exit(1)
			}
		}
	}
	fmt.Println("\033[1;32m✅ Helm is installed\033[0m")
}
