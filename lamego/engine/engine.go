package engine

import (
	"strings"

	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
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

func GenerateDefaultClass(m model.Model) target.Class {
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
			fieldType := ModelTypeToTargetType(f.GetTypeObject())
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
	}

	return c
	// What next?: Add variables | Add methods | Add getters/settings

}
