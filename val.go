package metla

import (
	"errors"
	_ "fmt"
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
	}
	return
}
