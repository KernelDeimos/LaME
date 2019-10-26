package engine

import (
	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/target"
)

type ModelReader interface {
	AddModel(m *model.Model)
}

type ModelProducer interface {
	SetModelReader(r ModelReader)
	Invoke()
}

type ClassReader interface {
	AddClass(c *target.Class)
}

type ClassProducer interface {
	SetClassReader(r ClassReader)
	Invoke()
}

type CodeProducer interface {
	SetAPI()
	Invoke()
}

type CodeProducerAPI interface {
	GetFileManager() target.FileManager
}

type EngineAPI interface {
	InstallModelReader(r ModelReader)
	InstallModelProducer(p ModelProducer)
	InstallClassReader(r ClassReader)
	InstallClassProducer(p ClassProducer)
	InstallCodeProducer(p CodeProducer)
}
