package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/virtualops/breeze-cli/pkg/attach"
	"github.com/virtualops/breeze-cli/pkg/docker"
	"github.com/virtualops/breeze-cli/pkg/kubernetes"
	"helm.sh/helm/v3/pkg/cli"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"strings"

	//"k8s.io/kubectl/pkg/cmd/attach"
	"k8s.io/kubectl/pkg/cmd/util"
	"os"
	"time"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an arbitrary command in your project",
	Long:  `The breeze run command creates a Kubernetes "pod" and attaches itself to the pod once running. Think of it like running a command over SSH.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must specify at least the command to run")
			os.Exit(1)
		}
		// Here, we build a docker setup slightly different to that of `dev`,
		// as we only want to spin up a single pod and attach to that.
		k8sClient, err := kubernetes.NewClient()
		if err != nil {
			fmt.Println("Failed to create Kubernetes client")
			os.Exit(1)
		}

		//appName := Config.ReleaseName()
		dockerClient, err := docker.NewClient()

		os.MkdirAll(".breezedev/run", 0777)
		f, err := os.Create(".breezedev/run/Dockerfile")

		if err != nil {
			fmt.Println("failed to create file", err)
		}

		// This only works for PHP based apps right now, and doesn't care about
		// the config at all. #startsmall
		dockerfile := &docker.Dockerfile{}
		dockerfile.Build(
			docker.PHP("7.4-fpm-alpine"),
			docker.ApkAdd("bash"),
			docker.Workdir("/var/www/html"),
			docker.Composer(),
			docker.Copy(".", "/var/www/html", docker.Chown("www-data", "www-data")),
			docker.ComposerAutoload,
		)

		f.WriteString(dockerfile.String())
		f.Close()

		imageTag, err := dockerClient.BuildImage(cmd.Context(), ".breezedev/run/Dockerfile", Config.ImageName())

		if err != nil {
			fmt.Println("Failed to build Docker image")
			os.Exit(1)
		}

		// here, we'll later do some check to see if the user wants to run the command in a cloud
		// environment, in which case we first need to push the image.

		// we probably also want a check here to make sure that we're *actually* in a local context
		// and otherwise do some magical switching logic

		// let's create a pod...
		pod, err := k8sClient.CoreV1().Pods("default").Create(cmd.Context(), &v1.Pod{
			Spec: v1.PodSpec{
				RestartPolicy: v1.RestartPolicyNever,
				Containers: []v1.Container{
					{
						Name:    "shell",
						Image:   imageTag,
						Command: []string{args[0]},
						Args:    args[1:],
						Stdin:   true,
						TTY:     true,
					},
				},
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-shell-%d", Config.ReleaseName(), time.Now().Unix()),
				Labels: map[string]string{
					"k8s.breeze.sh/type":     "shell",
					"app.kubernetes.io/name": Config.ReleaseName(),
				},
			},
		}, metav1.CreateOptions{
			FieldManager: "breeze.sh",
		})

		if err != nil {
			fmt.Println("failed to create pod")
			os.Exit(1)
		}

		s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
		s.Suffix = " Waiting for shell to boot"
		s.Start()

		labelSelectors := make([]string, 0)
		for k, v := range pod.Labels {
			labelSelectors = append(labelSelectors, fmt.Sprintf("%s=%s", k, v))
		}

		w, err := k8sClient.CoreV1().Pods("default").Watch(cmd.Context(), metav1.ListOptions{
			LabelSelector: strings.Join(labelSelectors, ","),
		})

		if err != nil {
			fmt.Println("Failed to watch for pods")
			os.Exit(1)
		}

		results := w.ResultChan()

		runCompleted := false
		for {
			event := <-results

			if event.Type != watch.Added && event.Type != watch.Modified {
				continue
			}
			updatedPod := event.Object.(*v1.Pod)

			if updatedPod.Status.Phase == v1.PodRunning {
				w.Stop()
				break
			}

			if updatedPod.Status.Phase == v1.PodSucceeded {
				w.Stop()
				runCompleted = true
				break
			}
		}

		// before attaching, we want to output the current logs from the pod, for any output that happened on boot
		res := k8sClient.CoreV1().Pods("default").GetLogs(pod.Name, &v1.PodLogOptions{
			Container: "shell",
		}).Do(cmd.Context())
		b, _ := res.Raw()
		fmt.Println(string(b))

		// if the command already completed without interaction, we will exit early
		if runCompleted {
			os.Exit(0)
		}

		// otherwise, we will try to attach to the container, assuming it's either an
		// interactive, or long-running process.
		factory := util.NewFactory(cli.New().RESTClientGetter())
		attachCmd := attach.NewCmdAttach(factory, genericclioptions.IOStreams{
			In:     os.Stdin,
			Out:    os.Stdout,
			ErrOut: os.Stderr,
		})

		// Run the attach command
		s.Stop()
		attachCmd.Flags().Parse([]string{"-i", "-t"})
		attachCmd.Run(attachCmd, []string{"pods", pod.Name})

		k8sClient.CoreV1().Pods("default").Delete(cmd.Context(), pod.Name, metav1.DeleteOptions{})
	},
}
