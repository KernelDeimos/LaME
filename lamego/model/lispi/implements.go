package lispi

/*
//::run : sequenceable (store (join-lf (DATA)))
func (i $1) AsSequenceableInstruction() SequenceableInstruction { return i }
//::end
//::run : expression (store (join-lf (DATA)))
func (i $1) AsExpressionInstruction() ExpressionInstruction { return i }
//::end
*/

//::gen repcsv (expression) StrLen StrSub StrCat
func (i StrLen) AsExpressionInstruction() ExpressionInstruction { return i }
func (i StrSub) AsExpressionInstruction() ExpressionInstruction { return i }
func (i StrCat) AsExpressionInstruction() ExpressionInstruction { return i }

//::end

//::gen repcsv (expression) Plus Minus Divide Multiply
func (i Plus) AsExpressionInstruction() ExpressionInstruction     { return i }
func (i Minus) AsExpressionInstruction() ExpressionInstruction    { return i }
func (i Divide) AsExpressionInstruction() ExpressionInstruction   { return i }
func (i Multiply) AsExpressionInstruction() ExpressionInstruction { return i }

//::end
