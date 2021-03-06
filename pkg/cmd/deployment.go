package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stefanprodan/kubectl-kubesec/pkg/kubesec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var deploymentCmd = &cobra.Command{
	Use:     `deployment [name]`,
	Short:   "Scans deployment object",
	Example: `  deployment podinfo --namespace=default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("deployment name is required")
		}
		name := args[0]

		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)

		fmt.Println("scanning deployment", name, "in namespace", namespace)
		dep, err := kubeClient.AppsV1beta2().Deployments(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dep.TypeMeta = metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		}
		err = serializer.Encode(dep, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		writer.Flush()

		result, err := kubesec.NewClient().ScanDefinition(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if result.Error != "" {
			fmt.Println(result.Error)
			os.Exit(1)
		}

		result.Dump(os.Stdout)

		return nil
	},
}
