package coreplugin

import (
	"strings"

	"github.com/KernelDeimos/LaME/lamego/engine"
	"github.com/KernelDeimos/LaME/lamego/engine/pluginapi"
	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

type DefaultClassGenerator struct {
	pluginapi.AbstractClassGenerator
	models      []model.Model
	classReader engine.ClassReader
}

func (p *DefaultClassGenerator) Install(api engine.EngineAPI) {
	api.InstallModelReader(p)
	api.InstallClassProducer(p)
	api.InstallUtilityUser(p)
}

func (p *DefaultClassGenerator) AddModel(m model.Model) {
	p.models = append(p.models, m)
}

func (p *DefaultClassGenerator) SetClassReader(r engine.ClassReader) {
	p.classReader = r
}

func (p *DefaultClassGenerator) InvokeClasses() []engine.ClassGenerationError {
	errorList := []engine.ClassGenerationError{}
	for _, m := range p.models {
		c, err := p.generateClass(m)
		if err != nil {
			errorList = append(errorList, *err)
			continue
		}
		p.classReader.AddClass(c)
	}
	p.models = []model.Model{}
	return errorList
}

func (p *DefaultClassGenerator) generateClass(m model.Model) (
	target.Class, *engine.ClassGenerationError,
) {
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
			fieldType := engine.ModelTypeToTargetType(f.GetTypeObject())
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
				Visibility: target.VisibilityPublic,
				Arguments:  []target.Variable{},
				Code: lispi.FakeBlock{
					StatementList: []lispi.SequenceableInstruction{
						lispi.Return{
							Expression: lispi.IGet{
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
				Visibility: target.VisibilityPublic,
				Arguments: []target.Variable{
					target.Variable{
						Name: "v",
						Type: fieldType,
					},
				},
				Code: lispi.FakeBlock{
					StatementList: []lispi.SequenceableInstruction{
						lispi.ISet{
							Name: issetName,
							Expression: lispi.LiteralBool{
								Value: true,
							},
						},
						lispi.ISet{
							Name: privateName,
							Expression: lispi.VGet{
								Name: "v",
							},
						},
					},
				},
			})
		}

		for _, me := range m.Methods {
			// TODO: throws panic on unrecognized syntax frontend
			sf := p.Utilities.SyntaxFrontends[m.Meta.GencodeSyntaxFrontend]

			var code string
			var codeSet bool

			// Prioritize language-targeted code
			/* not supported yet
			for lang, thisCode := range me.Hardcode {
				if lang == runConfig.TargetLanguage {
					code = thisCode
					codeSet = true
				}
			}
			*/

			var instructions []lispi.SequenceableInstruction
			if codeSet {
				instructions = []lispi.SequenceableInstruction{
					lispi.Raw{
						Value: code,
					},
				}
			} else {
				var err error
				instructions, err = sf.Process(me.Gencode)
				if err != nil {
					panic(err) // TODO
				}
			}

			methodArgs := []target.Variable{}
			for _, arg := range me.Args {
				methodArgs = append(methodArgs, target.Variable{
					Name: arg.Name,
					Type: engine.ModelTypeToTargetType(
						model.GetTypeObject(arg.Type),
					),
				})
			}

			c.Methods = append(c.Methods, target.Method{
				Name: me.Name,
				Visibility: engine.ModelVisibilityToTargetVisibility(
					me.Visibility),
				Return: target.Variable{
					Type: engine.ModelTypeToTargetType(
						model.GetTypeObject(me.Return),
					),
				},
				Arguments: methodArgs,
				Code: lispi.FakeBlock{
					StatementList: instructions,
				},
			})
		}

		if c.Meta.Serialize.JSON {
			out := lispi.StrCat{
				StringExpressionA: lispi.StrCat{
					StringExpressionA: lispi.LiteralString{Value: "{"},
					StringExpressionB: lispi.LiteralString{""},
				},
				StringExpressionB: lispi.LiteralString{Value: "}"},
			}
			doComma := false
			for _, f := range m.Fields {
				maybeComma := func(in lispi.ExpressionInstruction) lispi.ExpressionInstruction {
					if doComma {
						return lispi.StrCat{
							StringExpressionA: lispi.LiteralString{Value: ","},
							StringExpressionB: in,
						}
					}
					return in
				}
				t := model.GetTypeObject(f.Type)
				if t.Primitive != model.PrimitiveLaME {
					var expr lispi.ExpressionInstruction
					switch t.Primitive {
					case model.PrimitiveString:
						expr = lispi.StrCat{
							StringExpressionA: lispi.StrCat{
								StringExpressionA: lispi.LiteralString{Value: `"`},
								StringExpressionB: lispi.IGet{Name: f.Name + "__"},
							}, StringExpressionB: lispi.LiteralString{Value: `"`}}
					case model.PrimitiveInt:
						expr = lispi.IntToString{
							IntExpression: lispi.IGet{Name: f.Name + "__"},
						}
					}
					outA := out.StringExpressionA.(lispi.StrCat)
					outA.StringExpressionB = lispi.StrCat{
						StringExpressionA: outA.StringExpressionB,
						StringExpressionB: maybeComma(lispi.StrCat{
							StringExpressionA: lispi.LiteralString{Value: `"` + f.Name + `":`},
							StringExpressionB: expr})}
					out.StringExpressionA = outA
					doComma = true
				}
			}
			c.Methods = append(c.Methods, target.Method{
				Name: "serializeJSON",
				Return: target.Variable{
					Type: target.String,
				},
				Visibility: target.VisibilityPublic,
				Arguments:  []target.Variable{},
				Code: lispi.FakeBlock{
					StatementList: []lispi.SequenceableInstruction{
						lispi.Return{
							Expression: out,
						},
					},
				},
			})
		}

	}

	return c, nil
}
