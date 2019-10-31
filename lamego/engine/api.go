package engine

import (
	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/target"
)

type UtilityUser interface {
	SetUtilities(utilities UtilityPackage)
}

type UtilityPackage struct {
	SyntaxFrontends map[string]SyntaxFrontend
}

type RuntimeIntelligenceUser interface {
	SetRuntimeIntelligenceProvider(provider RuntimeIntelligenceProvider)
}

type RuntimeIntelligenceProvider interface {
	GetTypeMap(clsName, pkgName, methodName string) map[string]target.Type
}

type ModelReader interface {
	AddModel(m model.Model)
}

type ModelProducer interface {
	SetModelReader(r ModelReader)
	InvokeModels()
}

type ClassReader interface {
	AddClass(c target.Class)
}

type ClassProducer interface {
	SetClassReader(r ClassReader)
	InvokeClasses() []ClassGenerationError
}

type CodeProducer interface {
	SetAPI(api CodeProducerAPI)
	InvokeCodeGeneration() []CodeGenerationError
	EndFile(filename string, fm target.FileManager)
}

type CodeProducerAPI interface {
	GetFileManager() target.FileManager
}

type EngineAPI interface {
	InstallConfigurable(c Configurable)
	InstallUtilityUser(u UtilityUser)
	InstallModelReader(r ModelReader)
	InstallModelProducer(p ModelProducer)
	InstallClassReader(r ClassReader)
	InstallClassProducer(p ClassProducer)
	InstallCodeProducer(p CodeProducer)
	InstallRuntimeIntelligenceUser(u RuntimeIntelligenceUser)
}

type ConfigurationProvider interface {
	GetConfig(name string) string
}

type Configurable interface {
	SetConfig(provider ConfigurationProvider) *ConfigurationError
}

type EnginePlugin interface {
	Install(api EngineAPI)
}

type ErrorClass int

const (
	ErrorClassError ErrorClass = iota
	ErrorClassWarning
)

type ConfigurationError struct {
	M string
}

//::run : errtype (store model class code)
//::end

/*
//::run : errtmpl (store (join-lf (DATA)))
type $ucc-1GenerationError struct {
	Input$ucc-1ID string
	Output$ucc-1ID string
	InputLineNumber int
	OutputLineNumber int
	ErrorClass ErrorClass
}
//::end
*/

//::gen repcsv (errtmpl) (errtype)
type ModelGenerationError struct {
	InputModelID     string
	OutputModelID    string
	InputLineNumber  int
	OutputLineNumber int
	ErrorClass       ErrorClass
}
type ClassGenerationError struct {
	InputModelID     string
	InputClassID     string
	OutputClassID    string
	InputLineNumber  int
	OutputLineNumber int
	ErrorClass       ErrorClass
}
type CodeGenerationError struct {
	InputModelID     string
	InputClassID     string
	InputCodeID      string
	OutputCodeID     string
	InputLineNumber  int
	OutputLineNumber int
	ErrorClass       ErrorClass
}

//::end
