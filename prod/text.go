package prod

type text struct {
	position
	src []byte
}

func (s text) exec(exec *tplExec) *execError {
	return exec.Write(s.src)
}

func (s text) execType() execType {
	return execText
}

func (s text) String() string {
	return "{ text }"
}
