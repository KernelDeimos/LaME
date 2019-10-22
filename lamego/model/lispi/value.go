package lispi

type LiteralBool struct {
	Value bool
}

func (i LiteralBool) AsExpressionInstruction() ExpressionInstruction { return i }

type LiteralString struct {
	Value string
}

func (i LiteralString) AsExpressionInstruction() ExpressionInstruction { return i }
