package main

import (
	// "encoding/json"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/KernelDeimos/LaME/lamego/engine"
	// "github.com/KernelDeimos/LaME/lamego/generators"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/target"

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

	e := engine.NewEngine(config)
	e.ClassGenerators = map[string]target.ClassGenerator{
		"go": makeClassGeneratorGo(),
		"js": makeClassGeneratorJs(),
	}

	ee := e.RunAll()

	if ee != nil {
		logrus.Fatal(ee.String())
	}

}

func makeClassGeneratorJs() genjs.ClassGenerator {
	writeContext := support.NewWriteContext()

	return genjs.ClassGenerator{
		WriteContext: writeContext,
	}
}

func makeClassGeneratorGo() *gengo.ClassGenerator {
	writeContext := support.NewWriteContext()

	return &gengo.ClassGenerator{
		WriteContext: writeContext,
	}
}
