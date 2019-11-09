package lispi

type While struct {
	Condition ExpressionInstruction
	Code      SequenceableInstruction
}

func (i While) AsSequenceableInstruction() SequenceableInstruction { return i }

type Continue struct{}

func (i Continue) AsSequenceableInstruction() SequenceableInstruction { return i }

type Break struct{}

func (i Break) AsSequenceableInstruction() SequenceableInstruction { return i }

type If struct {
	Condition ExpressionInstruction
	Code      SequenceableInstruction
}

func (i If) AsSequenceableInstruction() SequenceableInstruction { return i }

func For(
	initialize SequenceableInstruction,
	condition ExpressionInstruction,
	postiterate SequenceableInstruction,
	code FakeBlock,
) SequenceableInstruction {
	code.StatementList = append(
		code.StatementList, postiterate)
	return FakeBlock{
		StatementList: []SequenceableInstruction{
			initialize,
			While{
				Condition: condition,
				Code:      code,
			},
		},
	}
}
