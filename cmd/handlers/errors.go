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

package handlers

import "errors"

var (
	ErrInvalidNumberOfArguments  = errors.New("you must provide the exact number of arguments required. Please refer to the help documentation for more details")
	ErrUnableToFindTemplate      = errors.New("stencil was unable to find the template. Please ensure the path is correct, and that the template exists at that path")
	ErrUnableToFindOutput        = errors.New("stencil was unable to find the output path")
	ErrMissingActionFlag         = errors.New("you must provide the --action flag when providing only one argument. Please refer to the help documentation for more details")
	ErrRepositoryNotFound        = errors.New("the repository specified by the --repo flag could not be found. Please ensure you have added the repository using 'stencil repo add' and that you have spelled the name correctly")
	ErrDuplicateTemplatePathArgs = errors.New("cannot specify template path when using the --action flag")
	ErrActionNotFound            = errors.New("the action specified by the --action flag could not be found. Please ensure you have defined the action in a .actions.stencil.yaml|yml|json file and that you have spelled the name correctly")
)
