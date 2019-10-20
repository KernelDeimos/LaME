package engine

import (
	"strings"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

type SyntaxFrontend interface{}

type Engine struct {
	Config EngineConfig
	SyntaxFrontends map[string]SyntaxFrontend
	ClassGenerators map[string]target.ClassGenerator
	TargetLanguage string
}

type EngineConfig struct {
	//
}

type EngineRunConfig struct {
	TargetLanguage string
	ModelSourceDirectory string
	GeneratorOutputDirectory string
}

func NewEngine(config EngineConfig) *Engine {
	e := Engine{
		Config: config,
	}
	return &e
}

type EngineError interface{
	String() string
}
type DeFactoEngineError struct{
	M string
}
func (ee DeFactoEngineError) String() string {
	return ee.M
}

func (e *Engine) Generate(runConfig EngineRunConfig) EngineError {
	allModels := []model.Model{}

	// Walk model source directory and load models
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
	err := filepath.Walk(
		runConfig.ModelSourceDirectory, callback)
	if err != nil {
		return DeFactoEngineError{M: err.Error()}
	}

	var fm target.FileManager = target.NewDeFactoFileManager(
		runConfig.GeneratorOutputDirectory,
		target.CursorConfig{
			IndentToken: "\t",
		},
	)

	// Get the expected ClassGenerator
	cg, exists := e.ClassGenerators[runConfig.TargetLanguage]
	if !exists {
		return DeFactoEngineError{M: "Unrecognized target language"}
	}

	for i := 0; i < len(allModels); i++ {
		m := allModels[i]
		c := e.GenerateDefaultClass(m)
		cg.WriteClass(c, fm)
	}

	for _, filename := range fm.GetFiles() {
		cg.EndFile(filename, fm)
	}

	fm.FlushAll()
	return nil
}

func (e *Engine) ModelTypeToTargetType(tIn model.Type) (tOut target.Type) {
	if tIn.Primitive == model.PrimitiveLaME {
		tOut.TypeOfType = target.LaMEType
		tOut.Identifier = tIn.Identifier
		return
	}

	tOut.TypeOfType = target.PrimitiveType

	m := map[model.Primitive]string{
		model.PrimitiveString: target.PrimitiveString,
		model.PrimitiveBool:   target.PrimitiveBool,
		model.PrimitiveInt:    target.PrimitiveInt,
		model.PrimitiveFloat:  target.PrimitiveFloat,
		model.PrimitiveObject: target.PrimitiveObject,
		model.PrimitiveVoid:   target.PrimitiveVoid,
	}

	tOut.Identifier = m[tIn.Primitive]
	return
}

func (e *Engine) GenerateDefaultClass(m model.Model) target.Class {
	// Get name & package from model ID
	c := target.Class{}
	{
		var pkgName, name string
		{
			idParts := strings.Split(m.ID, ".")
			l := len(idParts) - 1
			name = idParts[l]
			pkgName = strings.Join(idParts[:l], ".")
		}
		c.Name = name
		c.Package = pkgName
		c.Meta = target.ClassMeta{
			Serialize: target.SerializeMeta{
				JSON: true,
			},
		}

		c.Variables = []target.Variable{}
		c.Methods = []target.Method{}

		for _, f := range m.Fields {
			fieldType := e.ModelTypeToTargetType(f.GetTypeObject())
			privateName := f.Name + "__"
			issetName := privateName + "isSet"
			c.Variables = append(c.Variables, target.Variable{
				Name: privateName,
				Type: fieldType,
			})
			c.Variables = append(c.Variables, target.Variable{
				Name: issetName,
				Type: target.Bool,
			})
			c.Methods = append(c.Methods, target.Method{
				Name: "get" + util.String.Capitalize(f.Name),
				Return: target.Variable{
					Type: fieldType,
				},
				Arguments: []target.Variable{},
				Code: model.FakeBlock{
					StatementList: []model.SequenceableInstruction{
						model.Return{
							Expression: model.IGet{
								Name: privateName,
							},
						},
					},
				},
			})
			c.Methods = append(c.Methods, target.Method{
				Name: "set" + util.String.Capitalize(f.Name),
				Return: target.Variable{
					Type: target.Void,
				},
				Arguments: []target.Variable{
					target.Variable{
						Name: "v",
						Type: fieldType,
					},
				},
				Code: model.FakeBlock{
					StatementList: []model.SequenceableInstruction{
						model.ISet{
							Name: issetName,
							Expression: model.LiteralBool{
								Value: true,
							},
						},
						model.ISet{
							Name: privateName,
							Expression: model.VGet{
								Name: "v",
							},
						},
					},
				},
			})
		}

		/*
		for _, me := range m.Methods {
			// TODO: Need to know the target language
			//       here for choosing hardcode options
			// m.Meta.GencodeSyntaxFrontend
		}
		*/

		if c.Meta.Serialize.JSON {
			c.Methods = append(c.Methods, target.Method{
				Name: "serializeJSON",
				Return: target.Variable{
					Type: target.String,
				},
				Arguments: []target.Variable{},
				Code: model.FakeBlock{
					StatementList: []model.SequenceableInstruction{
						model.Return{
							Expression: model.ISerializeJSON{},
						},
					},
				},
			})
		}
	}

	return c
	// What next?: Add variables | Add methods | Add getters/settings

}
