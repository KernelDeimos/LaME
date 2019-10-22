package lispi

type ISerializeJSON struct{}

func (i ISerializeJSON) AsExpressionInstruction() ExpressionInstruction { return i }
