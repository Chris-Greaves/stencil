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
