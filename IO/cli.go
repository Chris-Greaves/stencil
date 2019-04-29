// Copyright © 2018 Christopher Greaves <cjgreaves97@hotmail.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package IO

import (
	"fmt"

	"github.com/Chris-Greaves/stencil/confighelper"
)

type CLI struct {
}

// GetOverrides will take all the settings and offer the user to override them
func (c CLI) GetOverrides(allSettings []confighelper.Setting) ([]confighelper.Setting, error) {
	var updatedSets []confighelper.Setting

	for _, setting := range allSettings {
		if output := offerSettingToUser(setting); output != "" {
			updatedSets = append(updatedSets, confighelper.Setting{Name: setting.Name, Value: output})
		}
	}

	return updatedSets, nil
}

func offerSettingToUser(setting confighelper.Setting) string {
	fmt.Printf("Conf Override: \"%v\" [%v]: ", setting.Name, setting.Value)
	output := ""
	fmt.Scanln(&output)
	return output
}
