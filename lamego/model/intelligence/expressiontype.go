package intelligence

import (
	"github.com/KernelDeimos/LaME/lamego/model/lispi"
	"github.com/KernelDeimos/LaME/lamego/target"
)

// GetTypeForExpression reports the type of a given LisPI expression, as well
// as any validation errors caused by invalid types of sub-expressions
// recursively.
// IF a class is provided, this function allows method-only LisPI nodes such as
//    lispi.ICall and lispi.IGet
// ELSE a method-only LisPI node results in an error.
func GetTypeForExpression(
	c *target.Class,
	vars *map[string]target.Type,
	errs *[]ExpressionError,
	ins lispi.ExpressionInstruction,
) target.Type {
	addErrorNotAMethod := func() {
		*errs = append(*errs, NewErrorNotAMethod())
	}
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
	case lispi.IGet:
		if c == nil {
			addErrorNotAMethod()
			return target.Void
		}
		for _, ivar := range c.Variables {
			if ivar.Name == specificIns.Name {
				return ivar.Type
			}
		}
		*errs = append(*errs, ExpressionError{
			ErrorID: ErrorUnrecognizedInstanceVariable,
			M:       "unrecognized instance variable '" + specificIns.Name + "'",
		})
		return target.Void
	case lispi.ICall:
		if c == nil {
			addErrorNotAMethod()
			return target.Void
		}
		for _, m := range c.Methods {
			if m.Name == specificIns.Name {
				return m.Return.Type
			}
		}
		*errs = append(*errs, ExpressionError{
			ErrorID: ErrorUnrecognizedMethod,
			M:       "unrecognized method '" + specificIns.Name + "'",
		})
		return target.Void
	case lispi.VGet:
		typ, exists := (*vars)[specificIns.Name]
		if !exists {
			*errs = append(*errs, ExpressionError{
				ErrorID: ErrorUnrecognizedVariable,
				M:       "unrecognized variable '" + specificIns.Name + "'",
			})
			return target.Void
		}
		return typ
	case lispi.ISerializeJSON:
		return target.String
	}
	// by default, generate void type for expression to produce error
	return target.Void
}
