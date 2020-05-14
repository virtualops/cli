package installer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"os/exec"
	"time"
)

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

const stableChartUrl = "https://kubernetes-charts.storage.googleapis.com"

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
	ensureStableRepoExists()
}

func ensureStableRepoExists() {
	var buf bytes.Buffer
	cmd := exec.Command("helm", "repo", "list", "-o", "json")
	cmd.Stdout = &buf
	repos := make([]*repositoryElement, 0)

	err := cmd.Run()

	if err == nil {
		if err := json.Unmarshal(buf.Bytes(), &repos); err != nil {
			fmt.Println("Failed to parse Helm repo list output")
			os.Exit(1)
		}
	}

	for _, repo := range repos {
		if repo.Name == "stable" {
			goto SUCCESS
		}
	}

	if err := exec.Command("helm", "repo", "add", "stable", stableChartUrl).Run(); err != nil {
		fmt.Println("\033[1;31m✘ Failed to add Helm's stable chart repository\033[0m")
	}

SUCCESS:
	fmt.Println("\033[1;32m✅ Helm's stable charts repository is configured\033[0m")
}
