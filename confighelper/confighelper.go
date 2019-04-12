package confighelper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
)

// Setting is a simple struct to represent a setting in a Conf. This is used when trying to Get or Set values in the Conf.
type Setting struct {
	Name  string
	Value string
}

// Conf encompasses the anonymous json object for a template config.
type Conf struct {
	raw *gabs.Container
}

// New will create a new Conf using the contents contain in the file found in "path"
func New(path string) (*Conf, error) {
	err := validateSettingsPath(path)
	if err != nil {
		return nil, fmt.Errorf("Error ocurred validating settings path. Error: %v", err.Error())
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error ocurred reading settings file. Error: %v", err.Error())
	}

	jsonParsed, err := gabs.ParseJSON(data)
	if err != nil {
		return nil, err
	}

	conf := Conf{raw: jsonParsed}

	return &conf, nil
}

// GetAllValues gets all of the overridable variables from the Conf
func (c *Conf) GetAllValues() ([]Setting, error) {
	children, err := c.raw.ChildrenMap()
	if err != nil {
		return nil, fmt.Errorf("Error ocurred getting root children from settings file. Error: %v", err.Error())
	}

	var sets []Setting

	getValuesOrCallChildren(children, &sets, "")

	return sets, nil
}

// SetValues will take an array of settings to put each one into the Conf. If the setting already exists it will update the value, else it will add the new setting.
func (c *Conf) SetValues(settings []Setting) error {
	var err error
	for _, setting := range settings {
		_, err = c.raw.SetP(setting.Value, setting.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// Object Returns the Conf as an anonymous object.
func (c *Conf) Object() interface{} {
	return c.raw.Data()
}

func validateSettingsPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return errors.New("Path to file does not exist")
	}
	if val := filepath.Ext(path); val != ".json" {
		return errors.New("Extension must be '.json'")
	}
	return nil
}

func getValuesOrCallChildren(children map[string]*gabs.Container, sets *[]Setting, objPath string) {
	for child := range children {
		nextChildren, _ := children[child].ChildrenMap()
		if len(nextChildren) < 1 {
			*sets = append(*sets, Setting{Name: objPath + child, Value: children[child].Data().(string)})
		} else {
			getValuesOrCallChildren(nextChildren, sets, objPath+child+".")
		}
	}
}
