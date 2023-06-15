/*
Copyright 2023 The Rook Authors. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"github.com/rook/kubectl-rook-ceph/pkg/k8sutil"
	"github.com/rook/kubectl-rook-ceph/pkg/logging"
	rookclient "github.com/rook/rook/pkg/client/clientset/versioned"
	"github.com/spf13/cobra"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	KubeConfig           string
	OperatorNamespace    string
	CephClusterNamespace string
)

// rookCmd represents the rook command
var RootCmd = &cobra.Command{
	Use:              "rook-ceph",
	Short:            "kubectl rook-ceph provides common management and troubleshooting tools for Ceph.",
	Args:             cobra.MinimumNArgs(1),
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// logging.Info("CephCluster namespace: %q", CephClusterNamespace)
		// logging.Info("Rook operator namespace: %q", OperatorNamespace)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {

	// Define your flags and configuration settings.

	RootCmd.PersistentFlags().StringVar(&KubeConfig, "kubeconfig", "", "kubernetes config path")
	RootCmd.PersistentFlags().StringVar(&OperatorNamespace, "operator-namespace", "rook-ceph", "Kubernetes namespace where rook operator is running")
	RootCmd.PersistentFlags().StringVarP(&CephClusterNamespace, "namespace", "n", "rook-ceph", "Kubernetes namespace where CephCluster is created")
}

func GetClientsets() *k8sutil.Clientsets {
	var err error

	clientsets := &k8sutil.Clientsets{}

	// 1. Create Kubernetes Client
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	clientsets.KubeConfig, err = kubeconfig.ClientConfig()
	if err != nil {
		logging.Fatal(err)
	}

	clientsets.Rook, err = rookclient.NewForConfig(clientsets.KubeConfig)
	if err != nil {
		logging.Fatal(err)
	}

	clientsets.Kube, err = k8s.NewForConfig(clientsets.KubeConfig)
	if err != nil {
		logging.Fatal(err)
	}

	return clientsets
}
