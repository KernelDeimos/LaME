package main

import (
	// "encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/KernelDeimos/LaME/lamego/engine"
	// "github.com/KernelDeimos/LaME/lamego/generators"
	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/target"

	"github.com/KernelDeimos/LaME/lamego/generators/gengo"
)

var simpleModelFile = `
- DEFINE_MODELS
- id: project.models.Passenger
  fields:
  - name
  - email
- id: project.models.Booking
  fields:
  - name: passenger
    type: project.models.Passenger
  - name: notes
    type: string
`

func main() {
	m := []model.Model{}
	b := []byte(simpleModelFile)
	err := yaml.Unmarshal(b, &m)
	if err != nil {
		logrus.Fatal(err)
	}

	for i := 0; i < len(m); i++ {
		fmt.Printf("=== Found model: %s\n", m[i].ID)
		c := engine.GenerateDefaultClass(m[i])
		/*
			dat, err := json.Marshal(c)
			if err != nil {
				logrus.Fatal(err)
			}
		*/
		// fmt.Println(string(dat))
		var codeCursor target.CodeCursor = target.NewStringCodeCursor("\t")
		cg := makeClassGeneratorGo()
		cg.WriteClass(c, codeCursor)
		fmt.Println(codeCursor.GetString())
	}
}

func makeClassGeneratorGo() gengo.ClassGenerator {
	writeContext := support.NewWriteContext()

	return gengo.ClassGenerator{
		WriteContext: writeContext,
	}
}
