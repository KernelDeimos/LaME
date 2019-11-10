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

// Note: two constants below are a work-in-progress JSON parser for LisPI,
//  but I decided to abandon this in favour of doing it after more advanced
//  syntax frontends are implemented.
const lispiSkipWhitespace = `
		(if (|| (|| (== (vget txtchr) (str '\n')) (== (vget txtchr) (str '\t'))) (== (vget txtchr (str ' ')))) (
			(vset i (+ (vget i) (int 1)))
			(continue)
		))
`

const LisPIDeserialize = `
(vset state (int 0))
(vset e (strlen (vget input)))
(vset i (int 0))
(while (< (vget i) (vget e)) (
	(vset txtchr (strsub input (vget i) (+ (vget i) (int 1))))
	(if (== (vget state) (int 0)) (` + /*state 0: wait for object*/ `
		` + lispiSkipWhitespace + `
		(if (! (== (vget txtchr) (str '{'))) (
			(vset errmsg (str 'Invalid character when expecting object start'))
			(return (vget errmsg))
		))
		(vset i (+ (vget i) (int 1)))
		(vset state (int 2))
		(continue)
	))
	(if (== (vget state) (int 2)) (` + /*state 1: wait for key*/ `
		` + lispiSkipWhitespace + `
		(if (! (== (vget txtchr) (str '"'))) (
			(vset errmsg (str 'Invalid character when expecting key start'))
			(return (vget errmsg))
		))
		(vset buffer (str ''))
		(vset escape (bool false))
		(vset i (+ (vget i) (int 1)))
		(vset state (int 3))
		(continue)
	))
	(if (== (vget state) (int 3)) (` + /*state 1: wait for key*/ `
		` + lispiSkipWhitespace + `
		(if (! (== (vget txtchr) (str '"'))) (
			(vset errmsg (str 'Invalid character when expecting key start'))
			(return (vget errmsg))
		))
		(vset keystart (vget i))
		(vset i (+ (vget i) (int 1)))
		(vset state (int 3))
		(continue)
	))
))

`

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
			SourceModel: m.ID,
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
			var defaultExpr lispi.ExpressionInstruction = lispi.IGet{
				Name: privateName,
			}
			if fieldType.TypeOfType == target.PrimitiveType {
				switch fieldType.Identifier {
				case target.PrimitiveString:
					defaultExpr = lispi.LiteralString{
						Value: "",
					}
				case target.PrimitiveInt:
					defaultExpr = lispi.LiteralInt{
						Value: "0", // TODO: why string?
					}
				case target.PrimitiveBool:
					defaultExpr = lispi.LiteralBool{
						Value: false,
					}
				}
			}
			c.Methods = append(c.Methods, target.Method{
				Name: "get" + util.String.Capitalize(f.Name),
				Return: target.Variable{
					Type: fieldType,
				},
				Visibility: target.VisibilityPublic,
				Arguments:  []target.Variable{},
				Code: lispi.FakeBlock{
					StatementList: []lispi.SequenceableInstruction{
						lispi.If{
							Condition: lispi.Not{
								A: lispi.IGet{Name: issetName},
							},
							Code: lispi.Return{
								Expression: defaultExpr,
							},
						},
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
			process := lispi.FakeBlock{
				StatementList: []lispi.SequenceableInstruction{},
			}
			out := lispi.StrCat{
				StringExpressionA: lispi.StrCat{
					StringExpressionA: lispi.LiteralString{Value: "{"},
					StringExpressionB: lispi.LiteralString{""},
				},
				StringExpressionB: lispi.LiteralString{Value: "}"},
			}
			doComma := false
			for _, f := range m.Fields {
				fieldGetter := lispi.ICall{
					Name: "get" + util.String.Capitalize(f.Name),
					Arguments: lispi.ExpressionList{
						Expressions: []lispi.ExpressionInstruction{},
					},
				}
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
								// TODO: quote escape
								StringExpressionB: fieldGetter,
							}, StringExpressionB: lispi.LiteralString{Value: `"`}}
					case model.PrimitiveInt:
						expr = lispi.IntToString{
							IntExpression: fieldGetter,
						}
					case model.PrimitiveBool:
						process.StatementList = append(
							process.StatementList,
							lispi.FakeBlock{
								StatementList: []lispi.SequenceableInstruction{
									lispi.VSet{
										Name: f.Name,
										Expression: lispi.LiteralString{
											Value: "false",
										},
									},
									lispi.If{
										Condition: fieldGetter,
										Code: lispi.VSet{
											Name: f.Name,
											Expression: lispi.LiteralString{
												Value: "true",
											},
										},
									},
								},
							})
						expr = lispi.VGet{Name: f.Name}
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
						process,
						lispi.Return{
							Expression: out,
						},
					},
				},
			})
			/*
				sf := parsing.SyntaxFrontendLisPINatural{}
				instructions, err := sf.Process(LisPIDeserialize)
				if err != nil {
					panic(err)
				}
				c.Methods = append(c.Methods, target.Method{
					Name: "deserializeJSON",
					Return: target.Variable{
						Type: target.String,
					},
					Visibility: target.VisibilityPublic,
					Arguments:  []target.Variable{},
					Code: lispi.FakeBlock{
						StatementList: instructions,
					},
				})
			*/
		}

	}

	return c, nil
}
