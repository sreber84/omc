/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package core

import (
	"fmt"
	"os"

	"github.com/gmeghnag/omc/cmd/helpers"
	"github.com/gmeghnag/omc/types"
	"github.com/gmeghnag/omc/vars"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/kubectl/pkg/describe"
	"sigs.k8s.io/yaml"
)

func describePod(currentContextPath string, defaultConfigNamespace string, resourceName string) {
	podPath := currentContextPath + "/namespaces/" + defaultConfigNamespace + "/pods/" + resourceName + "/" + resourceName + ".yaml"
	Pod := corev1.Pod{}
	_file := helpers.ReadYaml(podPath)
	if err := yaml.Unmarshal([]byte(_file), &Pod); err != nil {
		fmt.Fprintln(os.Stderr, "Error when trying to unmarshal file "+podPath)
		os.Exit(1)
	}
	fake := fake.NewSimpleClientset(&Pod)
	c := &types.DescribeClient{Namespace: defaultConfigNamespace, Interface: fake}
	d := describe.PodDescriber{c}
	out, _ := d.Describe(defaultConfigNamespace, resourceName, describe.DescriberSettings{ShowEvents: false})
	fmt.Printf(out)
}

var Pod = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"po", "pods"},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		resourceName := ""
		if len(args) == 1 {
			resourceName = args[0]
		}
		describePod(vars.MustGatherRootPath, vars.Namespace, resourceName)
	},
}
