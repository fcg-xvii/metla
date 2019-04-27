package prod

type execText struct {
	position
	src []byte
}

func (s execText) exec(exec *tplExec) *execError {
	return exec.Write(s.src)
}

func (s execText) String() string {
	return "{ text }"
}
