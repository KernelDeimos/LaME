package main

import (
	// "encoding/json"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/KernelDeimos/LaME/lamego/engine"
	// "github.com/KernelDeimos/LaME/lamego/generators"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/support"

	"github.com/KernelDeimos/LaME/lamego/generators/gengo"
)

func main() {
	var modelDir string
	var outputDir string
	var targetLanguage string
	{
		args := os.Args[1:]
		targetLanguage = args[0]
		modelDir = args[1]
		outputDir = args[2]
	}

	e := engine.NewEngine(engine.EngineConfig{
		//
	})
	e.ClassGenerators = map[string]target.ClassGenerator{
		"go": makeClassGeneratorGo(),
	}

	err := e.Generate(engine.EngineRunConfig{
		ModelSourceDirectory: modelDir,
		GeneratorOutputDirectory: outputDir,
		TargetLanguage: targetLanguage,
	})

	if err != nil {
		logrus.Fatal(err.String())
	}

}

func makeClassGeneratorGo() gengo.ClassGenerator {
	writeContext := support.NewWriteContext()

	return gengo.ClassGenerator{
		WriteContext: writeContext,
	}
}
