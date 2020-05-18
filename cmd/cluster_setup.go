package cmd

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage/driver"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"time"
)

var settings *cli.EnvSettings
var force bool
var mysql bool
var dbPort int
var helmConfig *action.Configuration

const breezeIngressReleaseName = "breeze-ingress"
const breezeMySQLReleaseName = "breeze-sql"

var clusterSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install base components into the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Check for existing Breeze ingress
		helmConfig = new(action.Configuration)
		settings = cli.New()
		helmDriver := "secret"
		if err := helmConfig.Init(settings.RESTClientGetter(), settings.Namespace(), helmDriver, func(format string, v ...interface{}) {}); err != nil {
			os.Exit(1)
		}

		installIngress()

		if mysql {
			installSQLDatabase()
		}
		os.Exit(0)
	},
}

func init() {
	clusterSetupCmd.Flags().BoolVarP(&force, "force", "f", false, "Force creation of a Breeze ingress, even if another ingress is detected")
	clusterSetupCmd.Flags().BoolVarP(&mysql, "mysql", "", true, "Run a mysql installation in the cluster")
	clusterSetupCmd.Flags().IntVarP(&dbPort, "db-port", "", 3306, "Specifies the local port to bind the database to")
}

func installIngress() {
	statusCommand := action.NewStatus(helmConfig)
	status, err := statusCommand.Run(breezeIngressReleaseName)
	// if there's a non-404 error, something went wrong and we'll exit out
	if err != nil && err != driver.ErrReleaseNotFound {
		fmt.Println("Failed to retrieve Helm release", err.Error())
		os.Exit(1)
	}

	// 1.1 If exists, print message, exit 0
	// This ignores the --force option
	if err == nil && status != nil {
		fmt.Println("\033[1;32m✅ Breeze ingress already configured\033[0m")
		return
	}

	// if we get here, the release was not found, so we'll check for other ingresses
	fmt.Println("No breeze ingress found, checking if another ingress is available")
	// 2. Check for other ingress-looking things (service with type=LoadBalancer)
	kubeClient, err := helmConfig.KubeClient.(*kube.Client).Factory.KubernetesClientSet()
	if err != nil {
		fmt.Println("Failed to connect to Kubernetes")
		os.Exit(1)
	}
	services, err := kubeClient.CoreV1().Services(settings.Namespace()).List(context.Background(), metaV1.ListOptions{})

	if err != nil {
		fmt.Println("Failed to list services from Kubernetes")
		os.Exit(1)
	}

	var existingIngress *v1.Service
	for _, svc := range services.Items {
		if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
			existingIngress = &svc
			break
		}
	}
	if existingIngress != nil {
		if force {
			// 2.1 if exists --force flag enabled, print warning, continue to 3
			fmt.Println(fmt.Sprintf("\u001B[1;33m⚠️  A possible ingress was found (%s), but --force has been enabled, so the Breeze ingress will be installed\033[0m", existingIngress.Name))
		} else {
			// 2.2 print message about existing, exit 1
			fmt.Println(fmt.Sprintf("\033[1;31m✘ A possible ingress was found (%s), the Breeze ingress will not be installed\033[0m", existingIngress.Name))
			os.Exit(1)
		}
	}
	// 3. install nginx-ingress with helm, release = breeze-ingress
	s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
	s.Suffix = " Installing Breeze ingress"
	s.Start()
	installCmd := action.NewInstall(helmConfig)
	installCmd.ReleaseName = breezeIngressReleaseName
	installCmd.Namespace = settings.Namespace()
	installCmd.Wait = true
	cp, err := installCmd.LocateChart("stable/nginx-ingress", settings)
	if err != nil {
		fmt.Println("Failed to locate ingress chart")
		os.Exit(1)
	}
	chartRequested, err := loader.Load(cp)
	if err != nil {
		fmt.Println("Failed to load chart")
		os.Exit(1)
	}
	_, err = installCmd.Run(chartRequested, map[string]interface{}{
		"controller": map[string]interface{}{
			"service": map[string]interface{}{
				"externalTrafficPolicy": "Local",
			},
		},
	})

	if err != nil {
		fmt.Println("Failed to install Helm chart", err.Error())
		os.Exit(1)
	}

	s.Stop()
	fmt.Println("\033[1;32m✅ Breeze ingress installed\033[0m")
}

func installSQLDatabase() {
	statusCommand := action.NewStatus(helmConfig)
	status, err := statusCommand.Run(breezeMySQLReleaseName)
	// if there's a non-404 error, something went wrong and we'll exit out
	if err != nil && err != driver.ErrReleaseNotFound {
		fmt.Println("Failed to retrieve Helm release", err.Error())
		os.Exit(1)
	}

	// 1.1 If exists, print message, exit 0
	if err == nil && status != nil {
		fmt.Println("\033[1;32m✅ Breeze DB already configured\033[0m")
		os.Exit(0)
	}

	// 3. install nginx-ingress with helm, release = breeze-sql
	s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
	s.Suffix = " Installing MySQL"
	s.Start()
	installCmd := action.NewInstall(helmConfig)
	installCmd.ReleaseName = breezeMySQLReleaseName
	installCmd.Namespace = settings.Namespace()
	installCmd.Wait = true
	cp, err := installCmd.LocateChart("stable/mysql", settings)
	if err != nil {
		fmt.Println("Failed to locate DB chart")
		os.Exit(1)
	}
	chartRequested, err := loader.Load(cp)
	if err != nil {
		fmt.Println("Failed to load chart")
		os.Exit(1)
	}

	serviceValues := map[string]interface{}{
		"type": "ClusterIP",
	}

	// if the user specifies `--db-port=0`, we leave the service as ClusterIP.
	// Otherwise, we will create a LoadBalancer with the port = dbPort. We
	// need to use a LoadBalancer over NodePort to use 127.0.0.1:3306.
	if dbPort != 0 {
		serviceValues["type"] = "LoadBalancer"
		serviceValues["port"] = dbPort
		serviceValues["loadBalancerIP"] = "localhost"
	}

	_, err = installCmd.Run(chartRequested, map[string]interface{}{
		"mysqlRootPassword": "root",
		"service":           serviceValues,
		"mysqlDatabase":     "breeze", // we'll create a default DB for convenience
		"mysqlUser":         "breeze",
		"mysqlPassword":     "breeze",
		"replication": map[string]interface{}{
			"enabled": false, // we don't need replication for our local env
		},
	})

	if err != nil {
		fmt.Println("Failed to install Helm chart", err.Error())
		os.Exit(1)
	}

	s.Stop()

	fmt.Println("\033[1;32m✅ Breeze DB installed\033[0m")
	if dbPort != 0 {
		fmt.Printf("You can reach your database on 127.0.0.1:%d\n", dbPort)
	} else {
		fmt.Println("Database was launched in cluster-only mode, and is only accessible by other applications in the cluster")
	}
}
