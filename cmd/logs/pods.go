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
package logs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gmeghnag/omc/cmd/helpers"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func logsPods(currentContextPath string, defaultConfigNamespace string, podName string, containerName string, previousFlag bool, allContainersFlag bool, logLevels []string) {
	var logFile string
	if previousFlag {
		logFile = "previous.log"
	} else {
		logFile = "current.log"
	}
	var _Items v1.PodList
	CurrentNamespacePath := currentContextPath + "/namespaces/" + defaultConfigNamespace
	_file, err := ioutil.ReadFile(CurrentNamespacePath + "/core/pods.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: namespace "+defaultConfigNamespace+" not found.")
		os.Exit(1)
	}
	if err := yaml.Unmarshal([]byte(_file), &_Items); err != nil {
		fmt.Fprintln(os.Stderr, "Error when trying to unmarshal file "+CurrentNamespacePath+"/core/pods.yaml")
		os.Exit(1)
	}
	podMatch := ""
	for _, Pod := range _Items.Items {
		if podName != Pod.Name {
			continue
		}
		podMatch = podName
		var containers []string
		containerMatch := ""
		if len(Pod.Spec.Containers) == 1 && containerName == "" {
			containerMatch = Pod.Spec.Containers[0].Name
		} else {
			if allContainersFlag {
				for _, c := range Pod.Spec.Containers {
					if len(logLevels) > 0 {
						FilterCatLogs(CurrentNamespacePath+"/pods/"+Pod.Name+"/"+c.Name+"/"+c.Name+"/logs/"+logFile, logLevels)
					} else {
						helpers.Cat(CurrentNamespacePath + "/pods/" + Pod.Name + "/" + c.Name + "/" + c.Name + "/logs/" + logFile)
					}
				}
				return
			} else {
				var containerSlice []v1.Container
				containerSlice = append(containerSlice, Pod.Spec.Containers...)
				containerSlice = append(containerSlice, Pod.Spec.InitContainers...)
				for _, c := range containerSlice {
					if containerName == c.Name {
						containerMatch = containerName
						break
					}
					containers = append(containers, c.Name)
				}
			}
		}
		if containerMatch == "" {
			for _, c := range Pod.Spec.InitContainers {
				if containerName == c.Name {
					containerMatch = containerName
					break
				}
			}
		}
		if containerMatch == "" {
			if containerName != "" {
				fmt.Fprintln(os.Stderr, "error: container "+containerName+" is not valid for pod "+Pod.Name)
				os.Exit(1)
			} else {
				fmt.Fprintln(os.Stderr, "error: a container name must be specified for pod "+Pod.Name+", choose one of:", containers)
				os.Exit(1)
			}
		} else {
			if len(logLevels) > 0 {
				FilterCatLogs(CurrentNamespacePath+"/pods/"+Pod.Name+"/"+containerMatch+"/"+containerMatch+"/logs/"+logFile, logLevels)
			} else {
				helpers.Cat(CurrentNamespacePath + "/pods/" + Pod.Name + "/" + containerMatch + "/" + containerMatch + "/logs/" + logFile)
			}
		}
	}
	if podMatch == "" {
		fmt.Fprintln(os.Stderr, "error: pods "+podName+" not found")
		os.Exit(1)
	}
}
