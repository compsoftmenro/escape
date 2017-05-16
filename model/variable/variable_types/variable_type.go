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

package variable_types

import (
	"errors"
)

type VariableType struct {
	Type            string
	UserCanOverride bool
	Validate        func(value interface{}, options map[string]interface{}) (interface{}, error)
}

func GetVariableType(typ string) (*VariableType, error) {
	knownTypes := map[string]func() *VariableType{
		"string":  NewStringVariableType,
		"integer": NewIntegerVariableType,
		"list":    NewListVariableType,
	}
	result, ok := knownTypes[typ]
	if !ok {
		return nil, errors.New("Variable type '" + typ + "' not implemented")
	}
	return result(), nil
}
