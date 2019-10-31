package lispi

type And struct {
	A ExpressionInstruction
	B ExpressionInstruction
}

type Or struct {
	A ExpressionInstruction
	B ExpressionInstruction
}

type Not struct {
	A ExpressionInstruction
}

type Lt struct {
	L ExpressionInstruction
	R ExpressionInstruction
}

// @avoid - syntax frontends should do (|| (< ..) (== ..))
type LtEq struct {
	L ExpressionInstruction
	R ExpressionInstruction
}

type Eq struct {
	A ExpressionInstruction
	B ExpressionInstruction
}

func (i And) AsExpressionInstruction() ExpressionInstruction  { return i }
func (i Or) AsExpressionInstruction() ExpressionInstruction   { return i }
func (i Not) AsExpressionInstruction() ExpressionInstruction  { return i }
func (i Lt) AsExpressionInstruction() ExpressionInstruction   { return i }
func (i LtEq) AsExpressionInstruction() ExpressionInstruction { return i }
func (i Eq) AsExpressionInstruction() ExpressionInstruction   { return i }

func Xor(A ExpressionInstruction, B ExpressionInstruction) ExpressionInstruction {
	return Or{
		A: And{
			A: A,
			B: Not{A: B},
		},
		B: And{
			A: Not{A: A},
			B: B,
		},
	}
}
