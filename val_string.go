/*
 *  Описание объекта строки и его креатора
 *
 */
package metla

import (
	"io"
	"reflect"
)

// Добавляем креатор строки в глобальный срез
/*func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValString,
		constructor: newValString,
	})
}*/

// Поверка соответствия строки - если открыта двойная или одинарная кавычка, тип соответствует
func checkValString(src []byte) bool {
	return src[0] == '\'' || src[0] == '"'
}

// Конструктор строки.
func newValString(p *parser) (res interface{}, err error) {
	p.SetupMark()
	charID := p.Char()     // Определяем, двойная или одинарная кавычка открыта
	p.IncPos()             // Смещаемся на начало строки
	if !p.ToChar(charID) { // Пытаемся найти закрывающую кавычку
		err = p.positionError("Unclosed string")
	} else {
		// Закрывающая кавычка найдена. Инициализируем результирующее значение
		p.stack.Push(p.MarkValString(0)[1:]) // Обрезаем кавычки для результирующего значения
		p.IncPos()
	}
	return
}

type valString struct {
	*rawInfoRecord
	val string
}

func (s *valString) Data(w io.Writer) (err error) {
	_, err = w.Write([]byte(s.val))
	return
}

func (s *valString) posInfo() *rawInfoRecord { return s.rawInfoRecord }
func (s *valString) String() string          { return "[string :: {" + s.val + "}]" }
func (s *valString) IsExecutable() bool      { return false }

func (s *valString) execObject(sto *storage, tpl *template, parent executor) (executor, error) {
	return s, nil
}

func (s *valString) Val() (interface{}, error) {
	return s.val, nil
}

func (s *valString) Vals() ([]interface{}, error) {
	return []interface{}{s.val}, nil
}

func (s *valString) receiveEvent(name string, params []interface{}) bool { return false }

func (s *valString) IsNil() bool        { return false }
func (s *valString) Type() reflect.Kind { return reflect.String }
func (s *valString) ValSingle() bool    { return true }
