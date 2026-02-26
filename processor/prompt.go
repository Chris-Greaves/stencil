// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package processor

import (
	"errors"

	"github.com/charmbracelet/huh"
)

func getFieldsForPrompt(p *Processor, fields []huh.Field) ([]huh.Field, error) {
	for key, prompt := range p.cfg.Vars.Prompt {
		switch prompt.Type {
		case "string":
			var value = ""
			if prompt.Default != "" {
				value = prompt.Default
			}
			fields = append(fields, promptUserForString(key, prompt))
			p.values[key] = value
		case "select(string)":
			var value = ""
			fields = append(fields, promptUserForStringSelect(key, prompt))
			p.values[key] = value
		default:
			return nil, errors.New("unsupported prompt type: " + prompt.Type)
		}
	}
	return fields, nil
}

func promptUserForString(key string, prompt ConfigPrompt) *huh.Input {
	input := huh.NewInput().
		Key(key).
		Title(key).
		Description(prompt.Description).
		Suggestions(prompt.Suggestions)

	if prompt.Default != "" {
		input.Accessor(NewDefaultAccessor(prompt.Default))
	}

	return input
}

func promptUserForStringSelect(key string, prompt ConfigPrompt) *huh.Select[string] {
	input := huh.NewSelect[string]().
		Title(key).
		Description(prompt.Description).
		Key(key)

	var options []huh.Option[string]
	for _, option := range prompt.Options {
		options = append(options, huh.NewOption(option, option))
	}

	return input.Options(options...)
}
