package genjs

import (
	"encoding/json"
	"strings"

	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

const constructorFunction = `
var constructor_ = function() {
	var obj = {};
	for (var i=0; i < this.fields.length; i++) {
		// TODO: null needs to be default value for type instead
		obj[this.fields[i].name] = null;
	}
	for (var i=0; i < this.methods.length; i++) {
		obj[this.methods[i].name] = this.methods[i].jsFunction;
	}
	return obj;
}
`

type ClassGenerator struct {
	WriteContext support.WriteContext
	Config       map[string]string

	objects map[string]struct{}
}

func (object ClassGenerator) SetConfig(config map[string]string) {
	object.Config = config
}

func (object ClassGenerator) WriteClass(
	c target.Class, fm target.FileManager,
) {
	filename := "generated_LaME.js"

	cc, isNew := fm.RequestFileForCode(filename)

	if isNew {
		object.objects = map[string]struct{}{}
		cc.AddLine("// GENERATED CODE - changes to this file may be overwritten")
		cc.AddLine("var project = {};")
		cc.AddLine(constructorFunction)
	}

	cname := c.Package + "." + c.Name
	cc.AddLine(cname + " = {};")

	cc.AddLine(cname + ".fields = [")
	cc.IncrIndent()
	for _, v := range c.Variables {
		// typ, isVoid := object.getTypeString(v.Type)
		cc.StartLine()
		object.writeVariable(cc, v)
		cc.AddString(",")
		cc.EndLine()
	}
	cc.DecrIndent()
	cc.AddLine("];")

	object.open(cc, cname+".methods = [")
	for _, m := range c.Methods {
		object.writeMethod(c, cc, m)
	}
	object.close(cc, "],", 1)

	cc.AddLine(cname + ".create = constructor_.bind(cname);")
}

func (object ClassGenerator) open(cc target.CodeCursor, txt string) {
	cc.AddLine(txt)
	cc.IncrIndent()
}

func (object ClassGenerator) close(cc target.CodeCursor, txt string, n int) {
	for i := 0; i < n; i++ {
		cc.DecrIndent()
		cc.AddLine(txt)
	}
}

func (object ClassGenerator) writeMethod(
	c target.Class, cc target.CodeCursor, m target.Method,
) {
	// typ, _ := object.getTypeString(m.Return.Type)

	object.open(cc, "{")
	cc.AddLine("name: " + `"` + m.Name + `",`)
	cc.StartLine()
	cc.AddString("typReturn: ")
	object.writeVariable(cc, m.Return)
	cc.AddString(",")
	cc.EndLine()
	argslice := []string{}
	for _, v := range m.Arguments {
		argslice = append(argslice, v.Name)
	}
	argString := strings.Join(argslice, ",")
	cc.AddLine("jsFunction: function (" + argString + ") {")
	cc.IncrIndent()
	object.writeFakeBlock(cc, m.Code)
	cc.DecrIndent()
	cc.AddLine("},")
	object.close(cc, "},", 1)
}

func (object ClassGenerator) writeVariable(cc target.CodeCursor, v target.Variable) {
	vv := Variable{
		Name: v.Name,
		Type: v.Type,
	}
	b, err := json.Marshal(vv)
	if err != nil {
		panic(err)
	}
	cc.AddString(string(b))
}

func (object ClassGenerator) EndFile(
	filename string, fm target.FileManager,
) {
	cc, _ := fm.RequestFileForCode(filename)
	cc.AddLine("module.exports = project;")
}

func (object ClassGenerator) getTypeString(t target.Type) (string, bool) {
	if t.TypeOfType == target.PrimitiveType &&
		t.Identifier == target.PrimitiveVoid {
		return "", true
	}
	return t.Identifier, false
}

type Variable struct {
	Name string      `json:"name"`
	Type target.Type `json:"type"`
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
	instance := "this"
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
		cc.AddString(";")
	case lispi.ISet:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString(instance + "." + specificIns.Name +
			" = ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	case lispi.If:
		cc.StartLine()
		cc.AddString("if ( ")
		object.writeExpressionInstruction(cc, specificIns.Condition)
		cc.AddString(" ) {")
		cc.EndLine()
		cc.IncrIndent()
		object.writeSequenceableInstruction(cc, specificIns.Code)
		cc.DecrIndent()
		cc.AddLine("}")
	case lispi.While:
		cc.StartLine()
		cc.AddString("for ( ")
		object.writeExpressionInstruction(cc, specificIns.Condition)
		cc.AddString(" ) {")
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
	switch specificIns := ins.(type) {
	case lispi.IGet:
		instance := "this"
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
		cc.AddString(" " + specificIns.Value + " ")
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
	case lispi.ISerializeJSON:
		// TODO: maybe replace this with a lispi statement
		// TODO: creating this anonymous function makes me
		//  wonder how expressions that require multiple
		//  statements of logic will be implemented for the
		//  target languages that don't support this.
		cc.AddString("JSON.stringify(this)")
	}
}
