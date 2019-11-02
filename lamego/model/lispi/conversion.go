package lispi

type IntToString struct {
	IntExpression ExpressionInstruction
}

func (i IntToString) AsExpressionInstruction() ExpressionInstruction { return i }
