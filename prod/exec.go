package prod

type coordinator interface {
	parseError(string) *parseError
	execError(string) *execError
}

type executer interface {
	coordinator
	exec(*tplExec) *execError
}

type getter interface {
	coordinator
	get(*tplExec) interface{}
}

type setter interface {
	coordinator
	set(*tplExec, interface{}) *execError
}

type position struct {
	tplName   string
	line, pos int
}

func (s position) parseError(text string) *parseError {
	return &parseError{s.tplName, s.line, s.pos, text}
}

func (s position) execError(text string) *execError {
	return &execError{s.tplName, s.line, s.pos, text}
}
