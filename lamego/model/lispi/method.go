package lispi

// TODO: IGetExpr will take an expression instead
type IGet struct {
	Name string
}

func (i IGet) AsExpressionInstruction() ExpressionInstruction { return i }

type ICall struct {
	Name      string
	Arguments ExpressionList
}

func (i ICall) AsExpressionInstruction() ExpressionInstruction     { return i }
func (i ICall) AsSequenceableInstruction() SequenceableInstruction { return i }

type ISet struct {
	Name       string
	Expression ExpressionInstruction
}

func (i ISet) AsSequenceableInstruction() SequenceableInstruction { return i }

type Get struct {
	Name string
}

func (i Get) AsExpressionInstruction() ExpressionInstruction { return i }

type Set struct {
	Name       string
	Expression ExpressionInstruction
}

func (i Set) AsSequenceableInstruction() SequenceableInstruction { return i }
