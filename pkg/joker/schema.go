package joker

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var _ Schema = (*JSchema)(nil)

type JSchema struct {
	Joker sJoker `yaml:"joker"`
}

type Schema interface {
	ParseSchema(path string) error
	PrepareUrlParams(p []map[any]any) string
}

type sJoker struct {
	Targets map[string]sTarget `yaml:"targets"`
}

type sTarget struct {
	Name   string        `yaml:"name"`
	Source sTargetSource `yaml:"source"`
	Id     string        `yaml:"id"`
	Read   []string      `yaml:"read"`
	Lang   sLang         `yaml:"lang"`
}

type sTargetSource struct {
	Method string        `yaml:"method"`
	Params []map[any]any `yaml:"params"`
	Target string        `yaml:"target"`
}

type sLang struct {
	Source    string `yaml:"source"`
	Target    string `yaml:"target"`
	Translate bool   `yaml:"translate"`
}

// Метод для создания строки http параметров
// из указанных в схеме
func (js *JSchema) PrepareUrlParams(p []map[any]any) string {
	var params string
	for k, v := range p {
		for keyName, keyValue := range v {
			if k < 1 {
				params = fmt.Sprintf("?%s=%s", keyName.(string), keyValue.(string))
				continue
			}

			params += fmt.Sprintf("&%s=%s", keyName.(string), keyValue.(string))
		}
	}

	return params
}

// Парсинг схемы джокера
func (js *JSchema) ParseSchema(path string) error {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yfile, js); err != nil {
		return err
	}

	// Сохранение названия целей
	for k := range js.Joker.Targets {
		if entry, ok := js.Joker.Targets[k]; ok {
			entry.Name = k
			js.Joker.Targets[k] = entry
		}
	}

	return nil
}
