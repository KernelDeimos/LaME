package engine

import (
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/target"
)

func (e *Engine) GenerateTypeMaps(c target.Class) []TypeValidationError {
	classID := c.Package + "." + c.Name
	errs := []TypeValidationError{}
	for _, m := range c.Methods {
		methodVars := map[string]target.Type{}
		e.genTypesForSequenceable(c, m, &methodVars, &errs, m.Code)
	}
}

func (e *Engine) genTypesForSequenceable(
	c target.Class, m target.Method,
	vars *map[string]target.Type,
	errs *[]TypeValidationError,
	ins lispi.SequenceableInstruction,
) {
	switch specificIns := ins.(type) {
	case lispi.Return:
		t := e.getTypeForExpression(
			c, m, vars, errs, specificIns.Expression)
		if t.TypeOfType != m.Return.Type.TypeOfType ||
			t.Identifier != m.Return.Type.Identifier {
			*errs = append(*errs, TypeValidationError{
				M: "return type mismatch",
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
	}
}

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
	case lispi.IGet:
		for _, ivar := range c.Variables {
			if ivar.Name == specificIns.Name {
				return ivar.Type
			}
		}
		*errs = append(*errs, TypeValidationError{
			M: "unrecognized instance variable '" + specificIns.Name + "'",
		})
		return target.Void
	case lispi.VGet:
		typ, exists := (*vars)[specificIns.Name]
		if !exists {
			*errs = append(*errs, TypeValidationError{
				M: "unrecognized variable '" + specificIns.Name + "'",
			})
			return target.Void
		}
	}
	// by default, generate void type for expression to produce error
	return target.Void
}
