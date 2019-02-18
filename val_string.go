/*
 *  Описание объекта строки и его креатора
 *
 */
package metla

import (
	"fmt"
)

// Добавляем креатор строки в глобальный срез
func init() {
	creators = append(creators, &valueCreator{
		checker:     checkValString,
		constructor: NewValString,
	})
}

// Поверка соответствия строки - если открыта двойная или одинарная кавычка, тип соответствует
func checkValString(src []byte) bool {
	return src[0] == '\'' || src[0] == '"'
}

// Конструктор строки.
func NewValString(p *parser) (res token, err error) {
	charID := p.char()     // Определяем, двойная или одинарная кавычка открыта
	p.incPos()             // Смещаемся на начало строки
	if !p.toChar(charID) { // Пытаемся найти закрывающую кавычку
		err = fmt.Errorf("Unclosed string (start position: [%v:%v])", p._mark.linePos, p._mark.pos)
	} else {
		// Закрывающая кавычка найдена. Инициализируем результирующее значение
		res = &valString{
			val: p.markValString(-1)[1:], // Обрезаем кавычки для результирующего значения
		}
		p.incPos()
	}
	return
}

type valString struct {
	val string
}

func (s *valString) Val() interface{} {
	return s.val
}

func (s *valString) Data() (res []byte, err error) {
	return []byte(s.val), nil
}

func (s *valString) String() string     { return "[string :: {" + s.val + "}]" }
func (s *valString) IsExecutable() bool { return false }
