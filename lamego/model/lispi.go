package model

// List Programming Interface

type SequenceableInstruction interface {
	AsSequenceableInstruction() SequenceableInstruction
}

type ExpressionInstruction interface {
	AsExpressionInstruction() ExpressionInstruction
}

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

type VGet struct {
	Name string
}

func (i VGet) AsExpressionInstruction() ExpressionInstruction { return i }

type LiteralBool struct {
	Value bool
}

func (i LiteralBool) AsExpressionInstruction() ExpressionInstruction { return i }

type LiteralString struct {
	Value string
}

func (i LiteralString) AsExpressionInstruction() ExpressionInstruction { return i }

type ISerializeJSON struct{}

func (i ISerializeJSON) AsExpressionInstruction() ExpressionInstruction { return i }

type Raw struct {
	Value string
}

func (i Raw) AsExpressionInstruction() ExpressionInstruction     { return i }
func (i Raw) AsSequenceableInstruction() SequenceableInstruction { return i }

/*
      return
        AND
          (<= this.min_length (lang.string.length args.value))
          (>= this.max_length (lang.string.length args.value))

lispi.Return(
	lispi.AND(
		lispi.LTEQ(
			lispi.IGET(
				"min_length",
				lispi.LCALL(
					"lang.string.length",
					lispi.GET("args.value"))))
		lispi.GTEQ(
			lispi.IGET(
				"max_length",
				lispi.LCALL(
					"lang.string.length",
					lispi.GET("args.value"))))
	)
)
*/
