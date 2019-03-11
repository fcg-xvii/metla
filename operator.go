package metla

// Типы операторов

type operatorType byte

const (
	opUndefined operatorType = iota
	opEndLine
	opArifmetic
	opFunction
	opSet
)

func checkOpType(ch byte) operatorType {
	switch ch {
	case '+', '-', '*', '/', '(':
		return opArifmetic
	case '=':
		return opSet
	case '\n':
		return opEndLine
	default:
		return opUndefined
	}
}

func isOpArifmetic(op byte) bool {
	return op == '+' || op == '-' || op == '*' || op == '/'
}
