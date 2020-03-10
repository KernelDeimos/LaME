package parsing

type Block struct {
	Statements []*Statement `(@@ ((","|";") @@)*)?`
}

type Statement struct {
	If     *IfStmt     `@@`
	While  *WhileStmt  `| @@`
	Set    *SetStmt    `| @@`
	Call   *CallStmt   `| @@`
	Return *ReturnStmt `| @@`
}

type IfStmt struct {
	Condition *Expression `"if " @@ ("("|"["|"{")`
	Code      *Block      `@@* (")"|"]"|"}")`
}

type WhileStmt struct {
	Condition *Expression `"while " @@ ("("|"["|"{")`
	Code      *Block      `@@* (")"|"]"|"}")`
}

type SetStmt struct {
	Name  string      `"=" @Ident`
	Value *Expression `@@`
}

type CallStmt struct {
	Name []string      `"x" @Ident ("." @Ident)*`
	Args []*Expression `(("("|"["|"{") (@@ ("," @@)*)? (")"|"]"|"}"))?`
}

type ReturnStmt struct {
	Value *Expression `"r" @@`
}

type Expression struct {
	String *string  `@String`
	Float  *string  `| @Float`
	Get    *GetExpr `| @@`
}

type GetExpr struct {
	Value string `@Ident`
}
