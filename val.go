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

func initCodeVal(p *parser) (err error) {
	fmt.Println("INIT_VAL", p.stack, string(p.Char()), string(p.EndLineContent()))
	p.PassSpaces()
	switch p.Char() {
	case '+', '-', '*', '/', '(', '!', '>', '<':
		err = newValArifmetic(p)
	case '"', '\'':
		err = newValString(p)
	case '=':
		err = newValSet(p)
	case '{':
		err = newValObject(p)
	case '[':
		err = newValArray(p)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err = newValNumber(p)
		fmt.Println("NUM", p.stack.Peek())
	default:
		if name, check := p.ReadName(); !check {
			err = p.positionError(fmt.Sprintf("Unexpected symbol '%c'", p.Char()))
		} else {
			//p.codeStack.Push(name)
			if keyword, check := getKeywordConstructor(string(name)); check {
				err = keyword(p)
			} else {
				switch p.Char() {
				case '(':
					err = newValFunction(p)
				case '[':
					err = newValIndex(p)
				case '.':
					err = newValField(p)
				default:
					p.stack.Push(&valVariable{p.infoRecordFromMark(), string(name)})
				}
			}
		}
	}
	return
}
