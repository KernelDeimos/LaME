package support

type ContextStackString []string

func (css *ContextStackString) Push(v string) {
	*css = append(*css, v)
}

func (css *ContextStackString) Unpush() {
	if len(*css) < 1 {
		panic("Tried to unpush empty ContextStackString")
	}
	*css = (*css)[:len(*css)-1]
}

func (css *ContextStackString) Get() string {
	return (*css)[len(*css)-1]
}

type WriteContext struct {
	ClassInstanceVariable ContextStackString
}

func NewWriteContext() WriteContext {
	return WriteContext{
		ClassInstanceVariable: ContextStackString{},
	}
}
