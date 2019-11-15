package intelligence

type ExpressionError struct {
	ErrorID ExpressionErrorID
	M       string
}

type ExpressionErrorID int

const (
	ErrorNULL ExpressionErrorID = iota
	ErrorUnrecognizedVariable
	ErrorUnrecognizedInstanceVariable
	ErrorUnrecognizedMethod
	ErrorNotAMethod
)

func NewErrorNotAMethod() ExpressionError {
	return ExpressionError{
		ErrorID: ErrorNotAMethod,
		M: "a non-method referred to an instance, " +
			"but non-methods do not have an instance.",
	}
}
