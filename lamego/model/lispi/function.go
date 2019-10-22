package lispi

// FakeBlock is a list of SequenceableInstruction thats contained
// inside some block (method, if statement, etc), but does not
// represent the block itself. For example, if multiple FakeBlock
// nodes appear in sequence, or within each other, they should be
// written to the same block of code.
type FakeBlock struct {
	StatementList []SequenceableInstruction
}

func (i FakeBlock) AsSequenceableInstruction() SequenceableInstruction { return i }

type Return struct {
	Expression ExpressionInstruction
}

func (i Return) AsSequenceableInstruction() SequenceableInstruction { return i }

type VGet struct {
	Name string
}

func (i VGet) AsExpressionInstruction() ExpressionInstruction { return i }
