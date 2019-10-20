package gengo

import (
	"fmt"
	"strings"

	"github.com/KernelDeimos/LaME/lamego/model"
	"github.com/KernelDeimos/LaME/lamego/support"
	"github.com/KernelDeimos/LaME/lamego/target"
	"github.com/KernelDeimos/LaME/lamego/util"
)

var MethodHeaderExpression = "func (o %s) %s(%s) %s"
var MethodHeaderVoid = "func (o %s) %s(%s)"

type ClassGenerator struct {
	WriteContext support.WriteContext
}

func (object ClassGenerator) asClass_lame_core_ClassGenerator() target.ClassGenerator {
	return object
}

func (object ClassGenerator) WriteClass(
	c target.Class, fm target.FileManager,
) {
	// Get code cursor
	cc := fm.RequestFileForCode(
		strings.Join(
			strings.Split(c.Package, "."),
			"/") + "/generated_LaME.go")

	// Add instance variable to the write context
	object.WriteContext.ClassInstanceVariable.Push("o")
	defer object.WriteContext.ClassInstanceVariable.Unpush()
	// Alright, writing a class in Go; here we go

	// Uhh.. I guess I gotta make a struct first
	cc.AddLine("type " + c.Name + " struct {")

	// Okay that worked, now variables
	func() {
		cc.IncrIndent()
		defer cc.DecrIndent()
		for _, v := range c.Variables {
			typ, isVoid := object.getTypeString(v.Type)
			if isVoid {
				// TODO: this is a user error
			}
			cc.AddLine(v.Name + " " + typ)
		}
	}()

	// Need to close the struct; methods go outside
	cc.AddLine("}")

	for _, m := range c.Methods {
		object.writeMethod(c, cc, m)
	}
}

func (object ClassGenerator) writeMethod(
	c target.Class, cc target.CodeCursor, m target.Method,
) {
	typ, isVoid := object.getTypeString(m.Return.Type)

	// Write the argument string
	var argString string
	{
		argslice := []string{}
		for _, v := range m.Arguments {
			argType, argTypeIsVoid :=
				object.getTypeString(v.Type)
			if argTypeIsVoid {
				// TODO: this is a user error
			}
			argslice = append(argslice,
				v.Name+" "+argType)
		}
		argString = strings.Join(argslice, ",")
	}

	if isVoid {
		cc.AddLine(fmt.Sprintf(MethodHeaderVoid,
			c.Name, m.Name, argString) + " {")
	} else {
		// Currently not supporting multiple return values
		// or named returns, since many target languages
		// won't support this anyway.
		returnString := typ
		cc.AddLine(fmt.Sprintf(MethodHeaderExpression,
			c.Name, m.Name, argString, returnString) + " {")
	}
	defer cc.AddLine("}")

	cc.IncrIndent()
	defer cc.DecrIndent()

	object.writeFakeBlock(cc, m.Code)
}

func (object ClassGenerator) getTypeString(t target.Type) (string, bool) {
	if t.TypeOfType == target.PrimitiveType &&
		t.Identifier == target.PrimitiveVoid {
		return "", true
	}
	return t.Identifier, false
}

func (object ClassGenerator) writeFakeBlock(
	cc target.CodeCursor, ins model.FakeBlock,
) {
	for _, subIns := range ins.StatementList {
		object.writeSequenceableInstruction(cc, subIns)
	}
}

func (object ClassGenerator) writeSequenceableInstruction(
	cc target.CodeCursor,
	ins model.SequenceableInstruction,
) {
	instance := object.WriteContext.ClassInstanceVariable.Get()
	// If this type switch thing raises any questions about
	// support for langauge definitions written in other
	// langauges, don't worry; I have a plan.
	switch specificIns := ins.(type) {
	case model.Return:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString("return ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	case model.ISet:
		cc.StartLine()
		defer cc.EndLine()
		cc.AddString(instance + "." + specificIns.Name +
			" = ")
		object.writeExpressionInstruction(cc, specificIns.Expression)
	}
}

func (object ClassGenerator) writeExpressionInstruction(
	cc target.CodeCursor,
	ins model.ExpressionInstruction,
) {
	switch specificIns := ins.(type) {
	case model.IGet:
		instance := object.WriteContext.ClassInstanceVariable.Get()
		cc.AddString(instance + ".get" +
			util.String.Capitalize(specificIns.Name))
	case model.VGet:
		cc.AddString(specificIns.Name)
	case model.LiteralBool:
		str := "false"
		if specificIns.Value {
			str = "true"
		}
		cc.AddString(str)
	}
}
