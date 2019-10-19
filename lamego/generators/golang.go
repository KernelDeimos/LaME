package generators

import (
	"github.com/KernelDeimos/LaME/lamego/core"
)

type ClassGenerator struct{}

func (object ClassGenerator) asClass_lame_core_ClassGenerator() core.ClassGenerator {
	return object
}

func (object ClassGenerator) WriteClass(m core.Model, cc core.CodeCursor) {

}
