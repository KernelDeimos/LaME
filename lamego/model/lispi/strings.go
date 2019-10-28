package lispi

type StrLen struct {
	StringExpression ExpressionInstruction
}

// Note: There is no CharAt. It is up to code generators to
//       optimize the equivalent special case of StrSub.
type StrSub struct {
	StringExpression ExpressionInstruction
	BeginAt          int
	EndAt            int
}

type StrCat struct {
	StringExpressionA ExpressionInstruction
	StringExpressionB ExpressionInstruction
}
