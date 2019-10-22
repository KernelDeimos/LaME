package lispi

// List Programming Interface

type SequenceableInstruction interface {
	AsSequenceableInstruction() SequenceableInstruction
}

type ExpressionInstruction interface {
	AsExpressionInstruction() ExpressionInstruction
}

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
