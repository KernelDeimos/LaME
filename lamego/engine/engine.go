package engine

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/parsing"
	"github.com/KernelDeimos/LaME/lamego/target"
)

func ModelTypeToTargetType(tIn model.Type) (tOut target.Type) {
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

type SyntaxFrontend interface {
	Process(script string) ([]lispi.SequenceableInstruction, error)
}

//::run : apis (store 'model reader' 'model producer' 'class reader' 'class producer' 'code producer')
//::end

//::run : readers (store model,model.Model class,target.Class)
//::end

type Engine struct {
	Config          EngineConfig
	SyntaxFrontends map[string]SyntaxFrontend
	ClassGenerators map[string]target.ClassGenerator
	//::gen repcsv '$ucc-1s []$ucc-1' (apis)
	ModelReaders   []ModelReader
	ModelProducers []ModelProducer
	ClassReaders   []ClassReader
	ClassProducers []ClassProducer
	CodeProducers  []CodeProducer
	//::end
	Configurables []Configurable
	UtilityUsers  []UtilityUser
	//::gen repcsv 'runtime$ucc-1List []$2' (readers)
	runtimeModelList []model.Model
	runtimeClassList []target.Class
	//::end
	// maps [class id][method name][variable name] -> target.Variable
	runtimeTypeMaps map[string]map[string]map[string]target.Type
	TargetLanguage  string

	fm target.FileManager
}

func (e *Engine) GetFileManager() target.FileManager {
	return e.fm
}

func (e *Engine) GetTypeMap(
	classID, methodName string) map[string]target.Type {
	// TODO: error handle
	return e.runtimeTypeMaps[classID][methodName]
}

type TypeValidationError struct {
	SourceClass  *target.Class
	SourceMethod *target.Method
	M            string
}

/*
//::run : api-setter (store (join-lf (DATA)))
func (e *Engine) Install$ucc-1($lcc-1 $ucc-1) {
	e.$ucc-1s = append(e.$ucc-1s, $lcc-1)
}
//::end
*/

//::gen repcsv (api-setter) (apis)
func (e *Engine) InstallModelReader(modelReader ModelReader) {
	e.ModelReaders = append(e.ModelReaders, modelReader)
}
func (e *Engine) InstallModelProducer(modelProducer ModelProducer) {
	e.ModelProducers = append(e.ModelProducers, modelProducer)
}
func (e *Engine) InstallClassReader(classReader ClassReader) {
	e.ClassReaders = append(e.ClassReaders, classReader)
}
func (e *Engine) InstallClassProducer(classProducer ClassProducer) {
	e.ClassProducers = append(e.ClassProducers, classProducer)
}
func (e *Engine) InstallCodeProducer(codeProducer CodeProducer) {
	e.CodeProducers = append(e.CodeProducers, codeProducer)
}

//::end

func (e *Engine) InstallConfigurable(configurable Configurable) {
	e.Configurables = append(e.Configurables, configurable)
}

func (e *Engine) InstallUtilityUser(user UtilityUser) {
	e.UtilityUsers = append(e.UtilityUsers, user)
}

/*
//::run : thing-adder (store (join-lf (DATA)))
func (e *Engine) Add$ucc-1($1 $2) {
	e.runtime$ucc-1List = append(e.runtime$ucc-1List, $1)
}
//::end
*/

//::gen repcsv (thing-adder) (readers)
func (e *Engine) AddModel(model model.Model) {
	e.runtimeModelList = append(e.runtimeModelList, model)
}
func (e *Engine) AddClass(class target.Class) {
	e.runtimeClassList = append(e.runtimeClassList, class)
}

//::end

type EngineConfig struct {
	Tasks []EngineRunConfig
}

type EngineDelegates struct {
	ClassGenerators map[string]target.ClassGenerator
}

type EngineRunConfig struct {
	Name                     string                            `yaml:"name"`
	TargetLanguage           string                            `yaml:"target"`
	ModelSourceDirectory     string                            `yaml:"source"`
	GeneratorOutputDirectory string                            `yaml:"output"`
	PluginConfig             map[string]map[string]interface{} `yaml:"config"`
}

func (conf EngineRunConfig) GetConfig(name string) string {
	b, err := json.Marshal(conf.PluginConfig[name])
	if err != nil {
		panic(err)
	}
	return string(b)
}

func NewEngine(config EngineConfig) *Engine {
	e := Engine{
		Config: config,
		SyntaxFrontends: map[string]SyntaxFrontend{
			"LisPI-Natural": parsing.SyntaxFrontendLisPINatural{},
		},
	}
	return &e
}

type EngineError interface {
	String() string
}
type DeFactoEngineError struct {
	M string
}

func (ee DeFactoEngineError) String() string {
	return ee.M
}

func (e *Engine) RunAll() EngineError {
	var err EngineError
	for _, conf := range e.Config.Tasks {
		fmt.Printf("[LaME] Running task: %s\n", conf.Name)
		err = e.Generate(conf)
		if err != nil {
			break
		}
	}
	return err
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

	e.fm = fm

	e.runtimeModelList = allModels

	// Configure plugins
	utilities := UtilityPackage{
		SyntaxFrontends: e.SyntaxFrontends,
	}
	for _, c := range e.Configurables {
		c.SetConfig(runConfig)
	}
	for _, u := range e.UtilityUsers {
		u.SetUtilities(utilities)
	}

	for _, p := range e.ModelProducers {
		p.SetModelReader(e)
	}
	for _, p := range e.ClassProducers {
		p.SetClassReader(e)
	}

	for i := 0; true; i++ {
		logrus.Infof("Invoking model producers (round %d)", i)
		for _, p := range e.ModelProducers {
			if p != nil {
				panic("todo")
			}
		}
		if len(e.runtimeModelList) < 1 {
			break
		}
		logrus.Infof("Feeding model readers (round %d)", i)
		for _, m := range e.runtimeModelList {
			for _, r := range e.ModelReaders {
				r.AddModel(m)
			}
		}
		e.runtimeModelList = []model.Model{}
	}

	e.runtimeClassList = []target.Class{}
	for i := 0; true; i++ {
		logrus.Infof("Invoking class producers (round %d)", i)
		for _, p := range e.ClassProducers {
			p.InvokeClasses()
		}
		if len(e.runtimeClassList) < 1 {
			break
		}
		// FINDME: update runtime type map
		for _, c := range e.runtimeClassList {
			errs := e.GenerateTypeMaps(c)
			if errs != nil && len(errs) > 0 {
				for _, thisErr := range errs {
					logrus.WithFields(logrus.Fields{
						"package": (*thisErr.SourceClass).Package,
						"class":   (*thisErr.SourceClass).Name,
						"method":  (*thisErr.SourceMethod).Name,
					}).Error(
						thisErr.M,
					)
				}
				logrus.Fatal("Halted with errors")
			}
		}
		logrus.Infof("Feeding class readers (round %d)", i)
		for _, c := range e.runtimeClassList {
			for _, r := range e.ClassReaders {
				r.AddClass(c)
			}
		}
		e.runtimeClassList = []target.Class{}
	}

	logrus.Infof("Invoking code producers")
	for _, p := range e.CodeProducers {
		p.SetAPI(e)
		p.InvokeCodeGeneration()
	}

	for _, filename := range fm.GetFiles() {
		for _, p := range e.CodeProducers {
			p.EndFile(filename, fm)
		}
	}

	fm.FlushAll()
	return nil
}
