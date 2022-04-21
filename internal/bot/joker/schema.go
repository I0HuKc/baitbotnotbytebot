package joker

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type JSchema struct {
	Joker sJoker `yaml:"joker"`
}

type Schema interface {
	ParseSchema(path string) error
	GetTarget(target string) sTarget
}

type sJoker struct {
	Targets map[string]sTarget `yaml:"targets"`
}

type sTarget struct {
	Name      string    `yaml:"name"`
	Target    string    `yaml:"target"`
	Id        sTargetId `yaml:"id"`
	Read      []string  `yaml:"read"`
	Translate bool      `yaml:"translate"`
}

type sTargetId struct {
	Field string `yaml:"field"`
	Type  string `yaml:"type"`
}

func (js *JSchema) GetTarget(target string) sTarget {
	return js.Joker.Targets[target]
}

func (js *JSchema) ParseSchema(path string) error {
	yfile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yfile, js)
}
