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

func checkOpType(lh, rh byte) operatorType {
	switch lh {
	case '+', '-', '*', '/', '(', '!':
		return opArifmetic
	case '=':
		if rh != '=' {
			return opSet
		}
		return opArifmetic
	case '\n':
		return opEndLine
	default:
		return opUndefined
	}
}
