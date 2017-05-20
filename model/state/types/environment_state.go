/*
Copyright 2017 Ankyra

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

package types

import (
	"fmt"
)

type EnvironmentState struct {
	Name        string                      `json:"name"`
	Inputs      map[string]interface{}      `json:"inputs"`
	Deployments map[string]*DeploymentState `json:"deployments"`
	ProjectName string                      `json:"-"`
}

func NewEnvironmentState(prjName, envName string) *EnvironmentState {
	return &EnvironmentState{
		ProjectName: prjName,
		Name:        envName,
		Inputs:      map[string]interface{}{},
		Deployments: map[string]*DeploymentState{},
	}
}

func (e *EnvironmentState) GetDeployments() []*DeploymentState {
	result := []*DeploymentState{}
	for _, d := range e.Deployments {
		result = append(result, d)
	}
	return result
}

func (e *EnvironmentState) GetProjectName() string {
	return e.ProjectName
}
func (e *EnvironmentState) getInputs() map[string]interface{} {
	return e.Inputs
}
func (e *EnvironmentState) GetName() string {
	return e.Name
}

func (e *EnvironmentState) ValidateAndFix(name, prjName string) error {
	e.Name = name
	e.ProjectName = prjName
	if e.Deployments == nil {
		e.Deployments = map[string]*DeploymentState{}
	}
	for deplName, depl := range e.Deployments {
		if err := depl.validateAndFix(deplName, e); err != nil {
			return err
		}
	}
	if e.ProjectName == "" {
		return fmt.Errorf("EnvironmentState's ProjectState reference has not been set")
	}
	if e.Name == "" {
		return fmt.Errorf("Environment name is missing from the EnvironmentState")
	}
	return nil
}

func (e *EnvironmentState) LookupDeploymentState(deploymentName string) (*DeploymentState, error) {
	val, ok := e.Deployments[deploymentName]
	if !ok {
		return nil, fmt.Errorf("Deployment '%s' does not exist", deploymentName)
	}
	return val, nil
}

func (e *EnvironmentState) GetDeploymentState(deps []string) (*DeploymentState, error) {
	if deps == nil || len(deps) == 0 {
		return nil, fmt.Errorf("Missing name to resolve deployment state. This is a bug in Escape.")
	}
	if len(deps) == 1 {
		return e.getOrCreateRootDeploymentState(deps[0]), nil
	} else {
		return e.getDeploymentStateForDependency(deps)
	}
}

func (e *EnvironmentState) getDeploymentStateForDependency(deps []string) (*DeploymentState, error) {
	deploymentName := deps[0]
	result := e.getOrCreateRootDeploymentState(deploymentName)
	for _, dep := range deps[1:] {
		depl, ok := result.Deployments[dep]
		if !ok {
			result = result.NewDependencyDeploymentState(dep)
		} else {
			result = depl
		}
	}
	return result, nil
}

func (e *EnvironmentState) getOrCreateRootDeploymentState(deploymentName string) *DeploymentState {
	depl, ok := e.Deployments[deploymentName]
	if !ok {
		depl = NewDeploymentState(e, deploymentName, deploymentName)
		e.Deployments[deploymentName] = depl
	}
	return depl
}