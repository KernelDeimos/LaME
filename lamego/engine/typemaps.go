package engine

import (
	"fmt"
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/target"
	// "github.com/sirupsen/logrus"
)

func (e *Engine) GenerateTypeMaps(c target.Class) []TypeValidationError {
	errs := []TypeValidationError{}
	for _, m := range c.Methods {
		methodVars := map[string]target.Type{}
		e.genTypesForSequenceable(c, m, &methodVars, &errs, m.Code)
		fmt.Println("--", c.Package, c.Name, m.Name)
		if _, isset := e.runtimeTypeMaps[c.Package+"."+c.Name]; !isset {
			e.runtimeTypeMaps[c.Package+"."+c.Name] =
				map[string]map[string]target.Type{}
		}
		e.runtimeTypeMaps[c.Package+"."+c.Name][m.Name] =
			methodVars
		fmt.Println("--", e.runtimeTypeMaps[c.Package+"."+c.Name][m.Name])
	}
	return errs
}

func (e *Engine) genTypesForSequenceable(
	c target.Class, m target.Method,
	vars *map[string]target.Type,
	errs *[]TypeValidationError,
	ins lispi.SequenceableInstruction,
) {
	switch specificIns := ins.(type) {
	case lispi.FakeBlock:
		for _, subIns := range specificIns.StatementList {
			e.genTypesForSequenceable(c, m, vars, errs, subIns)
		}
	case lispi.While:
		e.genTypesForSequenceable(c, m, vars, errs,
			specificIns.Code)
	case lispi.Return:
		t := e.getTypeForExpression(
			c, m, vars, errs, specificIns.Expression)
		if t.TypeOfType != m.Return.Type.TypeOfType ||
			t.Identifier != m.Return.Type.Identifier {
			*errs = append(*errs, TypeValidationError{
				M: "return type mismatch; " +
					"expected " + m.Return.Type.Identifier +
					" but got " + t.Identifier,
				SourceClass:  &c,
				SourceMethod: &m,
				// TODO: details
			})
		}
	case lispi.VSet:
		typ, varExists := (*vars)[specificIns.Name]
		t := e.getTypeForExpression(
			c, m, vars, errs, specificIns.Expression)
		if varExists {
			if t.TypeOfType != typ.TypeOfType ||
				t.Identifier != typ.Identifier {
				*errs = append(*errs, TypeValidationError{
					M: "variable type mismatch",
					// TODO: details
				})
			}
		} else {
			(*vars)[specificIns.Name] = t
		}
		// TODO: case lispi.Set:
	default:
		// logrus.Warn(ins)
	}
}

// @deprecated
// If this method is to be modified, use
// model/intelligence/GetTypeForExpression instead and
// then delete this method.
func (e *Engine) getTypeForExpression(
	c target.Class, m target.Method,
	vars *map[string]target.Type,
	errs *[]TypeValidationError,
	ins lispi.ExpressionInstruction,
) target.Type {
	switch specificIns := ins.(type) {
	case lispi.LiteralString:
		return target.String
	case lispi.LiteralBool:
		return target.Bool
	case lispi.LiteralInt:
		return target.Int
	case lispi.Plus:
		// TODO: search subtree for reals
		return target.Int
	case lispi.Minus:
		// TODO: search subtree for reals
		return target.Int
	case lispi.Multiply:
		// TODO: search subtree for reals
		return target.Int
	case lispi.Divide:
		// TODO: search subtree for reals
		return target.Int
	case lispi.StrLen:
		return target.Int
	case lispi.StrSub:
		return target.String
	case lispi.StrCat:
		return target.String
	case lispi.StrHash:
		return target.String
	case lispi.StrSalt:
		return target.String
	case lispi.IGet:
		for _, ivar := range c.Variables {
			if ivar.Name == specificIns.Name {
				return ivar.Type
			}
		}
		*errs = append(*errs, TypeValidationError{
			M:            "unrecognized instance variable '" + specificIns.Name + "'",
			SourceClass:  &c,
			SourceMethod: &m,
		})
		return target.Void
	case lispi.ICall:
		for _, m := range c.Methods {
			if m.Name == specificIns.Name {
				return m.Return.Type
			}
		}
		*errs = append(*errs, TypeValidationError{
			M:            "unrecognized method '" + specificIns.Name + "'",
			SourceClass:  &c,
			SourceMethod: &m,
		})
		return target.Void
	case lispi.VGet:
		typ, exists := (*vars)[specificIns.Name]
		if !exists {
			*errs = append(*errs, TypeValidationError{
				M:            "unrecognized variable '" + specificIns.Name + "'",
				SourceClass:  &c,
				SourceMethod: &m,
			})
			return target.Void
		}
		return typ
	case lispi.ISerializeJSON:
		return target.String
	case lispi.Eq:
		return target.Bool
	}
	// by default, generate void type for expression to produce error
	return target.Void
}
