package engine

import (
	"strings"

	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

func GenerateDefaultClass(m model.Model) {
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

		c.Variables = []target.Variable{}
		c.Methods = []target.Method{}

		for _, f := range m.Fields {
			c.Variables = append(c.Variables, target.Variable{
				Name: f.Name + "__",
				Type: f.GetTypeObject(),
			})
			c.Variables = append(c.Variables, target.Variable{
				Name: f.Name + "__isSet",
				Type: model.Bool,
			})
			c.Methods = append(c.Methods, target.Method{
				Name: "get" + util.String.Capitalize(f.Name),
			})
		}
	}

	// What next?: Add variables | Add methods | Add getters/settings

}
