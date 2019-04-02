package confighelper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Jeffail/gabs"
)

type Setting struct {
	Name  string
	Value string
}

type Conf struct {
	raw *gabs.Container
}

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

func (c *Conf) GetAllValues() ([]Setting, error) {
	children, err := c.raw.ChildrenMap()
	if err != nil {
		return nil, fmt.Errorf("Error ocurred getting root children from settings file. Error: %v", err.Error())
	}

	var sets []Setting

	getValuesOrCallChildren(children, &sets, "")

	return sets, nil
}

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

func (c *Conf) Object() interface{} {
	return c.raw.Data()
}

func (c *Conf) String() string {
	return c.raw.String()
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
			*sets = append(*sets, Setting{Name: fmt.Sprintf("%v%v", objPath, child), Value: children[child].Data().(string)})
		} else {
			getValuesOrCallChildren(nextChildren, sets, fmt.Sprintf("%v%v.", objPath, child))
		}
	}
}
