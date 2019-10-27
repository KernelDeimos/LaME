package main

import (
	// "encoding/json"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/KernelDeimos/LaME/lamego/engine"
	"github.com/KernelDeimos/LaME/lamego/engine/coreplugin"
	// "github.com/KernelDeimos/LaME/lamego/generators"
	"github.com/KernelDeimos/LaME/lamego/support"

	"github.com/KernelDeimos/LaME/lamego/generators/gengo"
	"github.com/KernelDeimos/LaME/lamego/generators/genjs"
)

func main() {
	var config engine.EngineConfig
	b, err := ioutil.ReadFile("LaME.yaml")
	if err != nil {
		logrus.Fatal(err)
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		logrus.Fatal(err)
	}

	for _, conf := range config.Tasks {
		e := engine.NewEngine(config)
		coreplugin.Plugin{}.Install(e)

		switch conf.TargetLanguage {
		case "go":
			makeClassGeneratorGo().Install(e)
		case "js":
			makeClassGeneratorJs().Install(e)
		}

		ee := e.Generate(conf)

		if ee != nil {
			logrus.Fatal(ee.String())
		}
	}
}

func makeClassGeneratorJs() *genjs.ClassGenerator {
	writeContext := support.NewWriteContext()

	return &genjs.ClassGenerator{
		WriteContext: writeContext,
	}
}

func makeClassGeneratorGo() *gengo.ClassGenerator {
	writeContext := support.NewWriteContext()

	return &gengo.ClassGenerator{
		WriteContext: writeContext,
	}
}
