/*
 *
 *
 */
package metla

// Конструктор строки.
func newValString(p *parser) (res interface{}, err error) {
	p.SetupMark()
	charID := p.Char()     // Определяем, двойная или одинарная кавычка открыта
	p.IncPos()             // Смещаемся на начало строки
	if !p.ToChar(charID) { // Пытаемся найти закрывающую кавычку
		err = p.positionError("Unclosed string")
	} else {
		// Закрывающая кавычка найдена. Инициализируем результирующее значение
		res = p.MarkValString(0)[1:]
		p.stack.Push(res) // Обрезаем кавычки для результирующего значения
		p.IncPos()
	}
	return
}
