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

func checkNumber(ch byte) bool {
	return ch >= 48 && ch <= 57
}

//Проверка соответствия первому символу наименования переменной (соответствует регулярке [0-9A-Za-z_])
func checkLetter(ch byte) bool {
	return (ch >= 97 && ch <= 122) || (ch >= 65 && ch <= 90) || ch == 95
}

// Шаблон checkLetter, на вхождеие добавлен символ '.'
func checkVarChar(ch byte) bool {
	return checkLetter(ch) || ch == '.'
}

////////////////////////////////////

// Идентификаторы операторов

////////////////////////////////////

type operator interface {
	Type() operatorType
	Data() ([]byte, error)
}
