package lispi

//::run : arithmeticops (store plus minus divide multiply)
//::end

/*
//::run : binaryop (store (join-lf (DATA)))
type $ucc-1 struct {
	A ExpressionInstruction
	B ExpressionInstruction
}
//::end
*/

//::gen repcsv (binaryop) (arithmeticops)
type Plus struct {
	A ExpressionInstruction
	B ExpressionInstruction
}
type Minus struct {
	A ExpressionInstruction
	B ExpressionInstruction
}
type Divide struct {
	A ExpressionInstruction
	B ExpressionInstruction
}
type Multiply struct {
	A ExpressionInstruction
	B ExpressionInstruction
}

//::end
