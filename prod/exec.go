package prod

type execType uint8

const (
	execFunction execType = iota
	execMethod
	execFor
	execBreak
	execIf
	execRPN
	execEcho
	execEcholn
	execPrint
	execText
	execField
	execSet
)

func (s execType) String() string {
	switch execFunction {
	case execFunction:
		return "function"
	case execMethod:
		return "method"
	case execFor:
		return "for"
	case execIf:
		return "if"
	case execRPN:
		return "rpn"
	case execEcho:
		return "echo"
	case execEcholn:
		return "echoln"
	case execPrint:
		return "print"
	case execText:
		return "text"
	case execField:
		return "field"
	default:
		return "undefined"
	}
}

type coordinator interface {
	parseError(string) *parseError
	execError(string) *execError
}

type executer interface {
	coordinator
	exec(*tplExec) *execError
	execType() execType
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
