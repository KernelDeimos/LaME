package lispi

type Raw struct {
	Value string
}

func (i Raw) AsExpressionInstruction() ExpressionInstruction     { return i }
func (i Raw) AsSequenceableInstruction() SequenceableInstruction { return i }
