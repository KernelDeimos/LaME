package lispi

// TODO: IGetExpr will take an expression instead
type IGet struct {
	Name string
}

func (i IGet) AsExpressionInstruction() ExpressionInstruction { return i }

type ISet struct {
	Name       string
	Expression ExpressionInstruction
}

func (i ISet) AsSequenceableInstruction() SequenceableInstruction { return i }
