package metla

import (
	"errors"
	"fmt"
)

type valueCheckMethod func([]byte) bool
type valueConstructor func(p *parser) (res token, err error)

type valueCreator struct {
	checker     valueCheckMethod
	constructor valueConstructor
}

var (
	creators               []*valueCreator
	errValueUnexpectedType = errors.New("Value init error :: Unexpected value type")
)

func getStartTypes(first []byte) (res []valueConstructor) {
	//fmt.Println(string(first))
	for _, creator := range creators {
		if creator.checker(first) {
			res = append(res, creator.constructor)
		}
	}
	return
}

func initVal(p *parser) (res token, err error) {
	/*
		//fmt.Println("INIT_VALL", string(p.EndLineContent()))
		// Если текущий символ соответствует завершению оператора или документа, это считается "пустым оператором". В даной ситуации ошибки не возникает
		p.PassSpaces()
		if p.IsEndLine() || p.IsEndDocument() || p.IsEndCode() {
			return
		}
		// Получаем данные от текущей позиции до конца строки и определяем возможные типы значений
		p.SetupMark()
		if types := getStartTypes(p.EndLineContent()); len(types) == 0 {
			err = p.positionError(errValueUnexpectedType.Error())
		} else {
			res, err = types[0](p)
		}*/
	return
}

func initCodeVal(p *parser) (val interface{}, err error) {
	fmt.Println("INIT_VAL", p.stack, string(p.Char()), string(p.EndLineContent()))
	p.PassSpaces()
	switch p.Char() {
	case '+', '-', '*', '/', '(', '!', '>', '<', '%':
		val, err = newValArifmetic(p)
	case '"', '\'':
		val, err = newValString(p)
	case '=':
		if p.NextChar() != '=' {
			val, err = newValSet(p)
		} else {
			val, err = newValArifmetic(p)
		}
	case '{':
		val, err = newValObject(p)
	case '[':
		val, err = newValArray(p)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, err = newValNumber(p)
	default:
		if name, check := p.ReadName(); !check {
			err = p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		} else {
			fmt.Println("NAME ACCEPTED")
			//p.codeStack.Push(name)
			if keyword, check := getKeywordConstructor(string(name)); check {
				val, err = keyword(p)
			} else {
				switch p.Char() {
				case '(':
					val, err = newValFunction(string(name), p)
				case '[':
					val, err = newValIndex(p)
				case '.':
					val, err = newValField(p)
				default:
					fmt.Println("VAL_VARIABLE")
					val = &valVariable{p.infoRecordFromMark(), string(name)}
					p.stack.Push(val)
				}
			}
		}
	}
	return
}
