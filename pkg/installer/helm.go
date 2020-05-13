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
		err = exec.Command("which", "brew").Run()

		fmt.Println("No Helm installation found, installing now")
		// use the official tilt installation, which prefers brew if available
		// and works with linux
		err = exec.Command("curl", "-fsSL", "https://raw.githubusercontent.com/windmilleng/tilt/master/scripts/install.sh", "-o", "install.sh").Run()
		if err != nil {
			fmt.Println("Failed to download Tilt installer", err)
			os.Exit(1)
		}
		s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = " Installing Tilt"
		s.Start()
		err = exec.Command("bash", "install.sh").Run()
		s.Stop()
		if err != nil {
			fmt.Println("Failed to install Tilt")
			os.Exit(1)
		}

		os.Remove(fmt.Sprintf("%s/%s", cwd, "install.sh"))
	}
	fmt.Println("\033[1;32m✅ Tilt is installed\033[0m")
}
