package gengo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KernelDeimos/LaME/lamego/engine"
	"github.com/KernelDeimos/LaME/lamego/engine/pluginapi"
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/support/golang"
	"github.com/KernelDeimos/LaME/lamego/support/typenamesupport"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

var MethodHeaderExpression = "func (o *%s) %s(%s) %s"
var MethodHeaderVoid = "func (o *%s) %s(%s)"

type FileState struct {
	imports map[string]struct{}
}

var fileStates map[string]*FileState

func init() {
	fileStates = map[string]*FileState{}
}

type ClassGenerator struct {
	pluginapi.AbstractCodeGenerator

	WriteContext support.WriteContext
	Config       map[string]interface{}
}

func (object *ClassGenerator) Install(api engine.EngineAPI) {
	api.InstallClassReader(object)
	api.InstallCodeProducer(object)
	api.InstallConfigurable(object)
	api.InstallRuntimeIntelligenceUser(object)
}

func (object *ClassGenerator) InvokeCodeGeneration() []engine.CodeGenerationError {
	for _, c := range object.ClassesToGenerate {
		object.WriteClass(c, object.CodeProducerAPI.GetFileManager())
	}
	return nil
}

func (object *ClassGenerator) SetConfig(
	provider engine.ConfigurationProvider,
) *engine.ConfigurationError {
	config := provider.GetConfig("language.go")
	err := json.Unmarshal([]byte(config), &object.Config)
	if err != nil {
		return &engine.ConfigurationError{M: err.Error()}
	}
	return nil
}

func (object ClassGenerator) WriteClass(
	c target.Class, fm target.FileManager,
) {
	filename := strings.Join(
		strings.Split(c.Package, "."),
		"/") + "/generated_LaME.go"

	// Get code cursor
	cc, isNew := fm.RequestFileForCode(filename)

	// Add instance variable to the write context
	object.WriteContext.ClassInstanceVariable.Push("o")
	defer object.WriteContext.ClassInstanceVariable.Unpush()
	// Alright, writing a class in Go; here we go

	var filestate *FileState
	if isNew {
		filestate = &FileState{}
		fileStates[filename] = filestate

		// Store imports for later writing to a subcursor
		filestate.imports = map[string]struct{}{}

		// Add package name and imports subcursor
		packageElems := strings.Split(c.Package, ".")
		packageName := packageElems[len(packageElems)-1]
		cc.AddLine("// GENERATED CODE - changes to this file may be overwritten")
		cc.AddLine("")
		cc.AddLine("package " + packageName)
		cc.NewSubCursor(support.CursorImports)
		cc.AddLine("")
	} else {
		filestate = fileStates[filename]
	}

	// Add implicit imports (from LaME core meta attributes)
	if c.Meta.Serialize.JSON {
		filestate.imports["encoding/json"] = struct{}{}
	}

	// Uhh.. I guess I gotta make a struct first
	cc.AddLine("type " + c.Name + " struct {")

	// Okay that worked, now variables
	func() {
		cc.IncrIndent()
		defer cc.DecrIndent()
		for _, v := range c.Variables {
			typ, isVoid := object.getTypeString(
				filestate, c, v.Type)
			if isVoid {
				// TODO: this is a user error
			}
			cc.AddLine(v.Name + " " + typ)
		}
	}()

	// Need to close the struct; methods go outside
	cc.AddLine("}")

	for _, m := range c.Methods {
		object.writeMethod(c, cc, m, filestate)
	}

}

func (object ClassGenerator) EndFile(
	filename string, fm target.FileManager,
) {
	// Get code cursor and file-related state
	cc, _ := fm.RequestFileForCode(filename)
	filestate := fileStates[filename]
	cc = cc.GetSubCursor(support.CursorImports)

	imports := filestate.imports

	// Write imports
	if len(imports) > 0 {
		cc.AddLine("")
		if len(imports) == 1 {
			for importString := range imports {
				cc.AddLine(`import "` + importString + `"`)
			}
		} else {
			func() {
				cc.AddLine(`import (`)
				defer cc.AddLine(`)`)
				cc.IncrIndent()
				defer cc.DecrIndent()
				for importString := range imports {
					cc.AddLine(`"` + importString + `"`)
				}
			}()
		}
	}
}

func (object ClassGenerator) writeMethod(
	c target.Class, cc target.CodeCursor, m target.Method,
	filestate *FileState,
) {
	typ, isVoid := object.getTypeString(filestate, c, m.Return.Type)

	mTargetName := m.Name
	if m.Visibility == target.VisibilityPublic {
		mTargetName = strings.Title(m.Name)
	}
	// used to skip arguments in variable declaration
	argNames := []string{}

	// Write the argument string
	var argString string
	{
		argslice := []string{}
		for _, v := range m.Arguments {
			argType, argTypeIsVoid :=
				object.getTypeString(filestate, c, v.Type)
			if argTypeIsVoid {
				// TODO: this is a user error
			}
			argslice = append(argslice,
				v.Name+" "+argType)
			argNames = append(argNames, v.Name)
		}
		argString = strings.Join(argslice, ",")
	}

	if isVoid {
		cc.AddLine(fmt.Sprintf(MethodHeaderVoid,
			c.Name, mTargetName, argString) + " {")
	} else {
		// Currently not supporting multiple return values
		// or named returns, since many target languages
		// won't support this anyway.
		returnString := typ
		cc.AddLine(fmt.Sprintf(MethodHeaderExpression,
			c.Name, mTargetName, argString, returnString) + " {")
	}
	defer cc.AddLine("}")

	// Write variable declarations
	vars := object.RuntimeIntelligenceProvider.GetTypeMap(
		c.Package, c.Name, m.Name)

	cc.IncrIndent()

	fmt.Println(c.Package, c.Name, m.Name)
	fmt.Println(vars)
	for name, typ := range vars {
		// TODO: typ.Identifier should not be used like this
		cc.AddLine("var " + name + " " + typ.Identifier)
	}
	cc.AddLine("")

	defer cc.DecrIndent()

	object.writeFakeBlock(cc, m.Code)
}

func (object ClassGenerator) getTypeString(filestate *FileState, c target.Class, t target.Type) (string, bool) {
	packages := golang.MapStrIToMapStrStr(
		object.Config["packages"].(map[string]interface{}))
	info := typenamesupport.GetTypeInfo(t, packages, c.Package)

	if info.IsPrimitive {
		if t.Identifier == target.PrimitiveVoid ||
			info.IsSpecialVoid {
			return "", true
		}
		return t.Identifier, false
	}

	if info.FailedToMatch {
		panic("Unrecognized package: " + info.LaMEPackage)
	}

	if info.IsCurrentPackage {
		return info.TypeName, false
	}

	filestate.imports[info.LanguagePackage] = struct{}{}
	var packageName string
	if len(info.LanguageRemainder) == 0 {
		parts := strings.Split(info.LanguagePackage, "/")
		packageName = parts[len(parts)-1]
	} else {
		packageName = info.LanguageRemainder[len(info.LanguageRemainder)-1]
	}

	return packageName + "." + info.TypeName, false
}

func (object ClassGenerator) writeFakeBlock(
	cc target.CodeCursor, ins lispi.FakeBlock,
) {
	for _, subIns := range ins.StatementList {
		object.writeSequenceableInstruction(cc, subIns)
	}
}

func (object ClassGenerator) writeSequenceableInstruction(
	cc target.CodeCursor,
	ins lispi.SequenceableInstruction,
) {
	instance := object.WriteContext.ClassInstanceVariable.Get()
	// If this type switch thing raises any questions about
	// support for langauge definitions written in other
	// langauges, don't worry; I have a plan.
	switch specificIns := ins.(type) {
	case lispi.FakeBlock:
		for _, ins := range specificIns.StatementList {
			object.writeSequenceableInstruction(cc, ins)
		}
	case lispi.Return:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString("return ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	case lispi.ISet:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString(instance + "." + specificIns.Name +
			" = ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	case lispi.VSet:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString(specificIns.Name + " = ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	case lispi.If:
		cc.StartLine()
		cc.AddString("if ")
		object.writeExpressionInstruction(cc, specificIns.Condition)
		cc.AddString(" {")
		cc.EndLine()
		cc.IncrIndent()
		object.writeSequenceableInstruction(cc, specificIns.Code)
		cc.DecrIndent()
		cc.AddLine("}")
	case lispi.While:
		cc.StartLine()
		cc.AddString("for ")
		object.writeExpressionInstruction(cc, specificIns.Condition)
		cc.AddString(" {")
		cc.EndLine()
		cc.IncrIndent()
		object.writeSequenceableInstruction(cc, specificIns.Code)
		cc.DecrIndent()
		cc.AddLine("}")
	}
}

func (object ClassGenerator) writeExpressionInstruction(
	cc target.CodeCursor,
	ins lispi.ExpressionInstruction,
) {
	/*
		//::run : infixu (store (join-lf (DATA)))
		case lispi.$1:
			object.writeExpressionInstruction(cc, specificIns.A)
			cc.AddString(" $2 ")
			object.writeExpressionInstruction(cc, specificIns.B)
		//::end
		//::run : infixu-ops (store (DATA))
		Plus,+
		Minus,-
		Divide,/
		Multiply,*
		//::end
	*/
	switch specificIns := ins.(type) {
	case lispi.IGet:
		instance := object.WriteContext.ClassInstanceVariable.Get()
		cc.AddString(instance + ".get" +
			util.String.Capitalize(specificIns.Name) + "()")
	case lispi.VGet:
		cc.AddString(specificIns.Name)
	case lispi.LiteralBool:
		str := "false"
		if specificIns.Value {
			str = "true"
		}
		cc.AddString(str)
	case lispi.LiteralInt:
		cc.AddString(specificIns.Value)
	case lispi.And:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" && ")
		object.writeExpressionInstruction(cc, specificIns.B)
	case lispi.Or:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" || ")
		object.writeExpressionInstruction(cc, specificIns.B)
	case lispi.Lt:
		// Note: model validator will eventually ensure
		//       the arguments are always integer types
		object.writeExpressionInstruction(cc, specificIns.L)
		cc.AddString(" < ")
		object.writeExpressionInstruction(cc, specificIns.R)
	case lispi.LtEq:
		// Note: model validator will eventually ensure
		//       the arguments are always integer types
		object.writeExpressionInstruction(cc, specificIns.L)
		cc.AddString(" <= ")
		object.writeExpressionInstruction(cc, specificIns.R)
	case lispi.Not:
		cc.AddString("!(")
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(")")
	case lispi.Eq:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" == ")
		object.writeExpressionInstruction(cc, specificIns.B)
	//::gen repcsv (infixu) (infixu-ops)
	case lispi.Plus:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" + ")
		object.writeExpressionInstruction(cc, specificIns.B)
	case lispi.Minus:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" - ")
		object.writeExpressionInstruction(cc, specificIns.B)
	case lispi.Divide:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" / ")
		object.writeExpressionInstruction(cc, specificIns.B)
	case lispi.Multiply:
		object.writeExpressionInstruction(cc, specificIns.A)
		cc.AddString(" * ")
		object.writeExpressionInstruction(cc, specificIns.B)
	//::end
	case lispi.StrSub:
		cc.AddString("(")
		object.writeExpressionInstruction(cc, specificIns.StringExpression)
		cc.AddString(")[(")
		object.writeExpressionInstruction(cc, specificIns.BeginAt)
		cc.AddString("):(")
		object.writeExpressionInstruction(cc, specificIns.EndAt)
		cc.AddString(")]")
	case lispi.StrLen:
		cc.AddString("len(")
		object.writeExpressionInstruction(cc, specificIns.StringExpression)
		cc.AddString(")")
	case lispi.ISerializeJSON:
		// TODO: maybe replace this with a lispi statement
		// TODO: creating this anonymous function makes me
		//  wonder how expressions that require multiple
		//  statements of logic will be implemented for the
		//  target languages that don't support this.
		cc.AddString("(func() string {")
		cc.EndLine()
		cc.IncrIndent()
		cc.AddLine("bout, err := json.Marshal(o)")
		cc.AddLine("if err != nil { return \"\" }")
		cc.AddLine("return string(bout)")
		cc.DecrIndent()
		cc.StartLine()
		cc.AddString("})()")
	}
}
