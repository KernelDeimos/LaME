package coreplugin

import (
	"github.com/KernelDeimos/LaME/lamego/engine"
	"github.com/KernelDeimos/LaME/lamego/model"
)

type Plugin struct{}

func (p Plugin) Install(api engine.EngineAPI) {
	a := DefaultClassGenerator{
		models: []model.Model{},
	}

	a.Install(api)
}
