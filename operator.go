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
	case '+', '-', '*', '/', '>', '<', '!':
		return opArifmetic
	case '(':
		return opFunction
	case '=':
		return opSet
	case '\n':
		return opEndLine
	default:
		return opUndefined
	}
}

////////////////////////////////////

// Идентификаторы операторов

////////////////////////////////////

type operator interface {
	Type() operatorType
	Data() ([]byte, error)
}
