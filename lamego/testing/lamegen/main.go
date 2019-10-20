package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/KernelDeimos/LaME/lamego/engine"
	// "github.com/KernelDeimos/LaME/lamego/generators"
	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/target"

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

	allModels := []model.Model{}

	callback := func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			m := []model.Model{}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(b, &m)
			if err != nil {
				return err
			}
			allModels = append(allModels, m[1:]...)
		}
		return nil
	}
	err := filepath.Walk(modelDir, callback)
	if err != nil {
		logrus.Fatal(err)
	}

	var fm target.FileManager = target.DeFactoFileManager{
		RootPath: outputDir,
		CursorConfig: target.CursorConfig{
			IndentToken: "\t",
		},
	}
	for i := 0; i < len(allModels); i++ {
		m := allModels[i]
		fmt.Printf("=== Found model: %s\n", m.ID)
		c := engine.GenerateDefaultClass(m)
		var cg target.ClassGenerator
		switch targetLanguage {
		case "go":
			cg = makeClassGeneratorGo()
		}
		cg.WriteClass(c, fm)
	}
}

func makeClassGeneratorGo() gengo.ClassGenerator {
	writeContext := support.NewWriteContext()

	return gengo.ClassGenerator{
		WriteContext: writeContext,
	}
}
