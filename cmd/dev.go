package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/virtualops/breeze-cli/pkg/config"
	"github.com/virtualops/breeze-cli/pkg/installer"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var devCmd = &cobra.Command{
	Use: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		installer.InstallOrVerifyTilt()
		installer.InstallOrVerifyHelm()
		buildDevFiles()
		buildK8sFiles()
		go removeDevFilesOnExit()
		RunTilt()
	},
}

func buildDevFiles() {
	cwd, _ := os.Getwd()
	dir := filepath.Base(cwd)
	releaseName := fmt.Sprintf("%s-%s", Config.Project, dir)
	dockerImageName := fmt.Sprintf("%s/%s", Config.Project, dir)
	os.Mkdir(".breezedev", 0777)
	f, err := os.Create(".breezedev/Dockerfile")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Dockerfile)
	f.Close()
	f, err = os.Create(".breezedev/nginx.Dockerfile")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.NginxDockerfile)
	f.Close()

	f, err = os.Create(".breezedev/opcache.ini")
	f.WriteString(config.Opcache)
	f.Close()

	f, err = os.Create(".breezedev/preload.php")
	f.WriteString(config.PreloadClass)
	f.Close()

	f, err = os.Create(".breezedev/Tiltfile")
	f.WriteString(`k8s_yaml(helm(
  'kubernetes',
  name='` + releaseName + `',
  set=['image=` + dockerImageName + `', 'nginxImage=` + dockerImageName + `-nginx', 'service.path=` + Config.Deploy.Path + `']
))
docker_build('` + dockerImageName + `', '..', dockerfile='Dockerfile',ignore=['/vendor'])
docker_build('` + dockerImageName + `-nginx', '../public', dockerfile='nginx.Dockerfile')`)
	f.Close()

	fmt.Println("Generated dev files")
}

func buildK8sFiles() {
	os.MkdirAll(".breezedev/kubernetes/templates", 0777)
	f, err := os.Create(".breezedev/kubernetes/templates/deployment.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Deployment)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/templates/service.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Service)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/templates/secret.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Secret)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/templates/nginx-config.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.NginxConfig)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/templates/fpm-config.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.FPMConfigMap)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/templates/ingress.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Ingress)
	f.Close()

	f, err = os.Create(".breezedev/kubernetes/Chart.yaml")

	if err != nil {
		fmt.Println("failed to create file", err)
	}
	f.WriteString(config.Chart)
	f.Close()
}

func RunTilt() {
	cmd := exec.Command("tilt", "up", "-f", ".breezedev/Tiltfile")
	cmd.Stdout = io.MultiWriter(os.Stdout)
	err := cmd.Run()

	if err != nil {
		fmt.Println("Tilt failed")
	}
}

func removeDevFilesOnExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	<-c
	err := os.RemoveAll(".breezedev")
	if err != nil {
		fmt.Println(err)
	}
}
