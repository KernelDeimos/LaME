package pluginapi

import (
	"github.com/KernelDeimos/LaME/lamego/engine"
	"github.com/KernelDeimos/LaME/lamego/target"
)

type AbstractCodeGenerator struct {
	ClassesToGenerate           []target.Class
	CodeProducerAPI             engine.CodeProducerAPI
	RuntimeIntelligenceProvider engine.RuntimeIntelligenceProvider
}

func (plugin *AbstractCodeGenerator) AddClass(c target.Class) {
	plugin.ClassesToGenerate = append(plugin.ClassesToGenerate, c)
}

func (plugin *AbstractCodeGenerator) SetAPI(api engine.CodeProducerAPI) {
	plugin.CodeProducerAPI = api
}

func (plugin *AbstractCodeGenerator) SetRuntimeIntelligenceProvider(
	provider engine.RuntimeIntelligenceProvider,
) {
	plugin.RuntimeIntelligenceProvider = provider
}

type AbstractClassGenerator struct {
	Utilities engine.UtilityPackage
}

func (plugin *AbstractClassGenerator) SetUtilities(utilities engine.UtilityPackage) {
	plugin.Utilities = utilities
}
