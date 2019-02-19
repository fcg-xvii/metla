package metla

import (
	"errors"
	"fmt"
)

var (
	parseErrEOF            = errors.New("EOF")
	parseUnexpectedLiteral = errors.New("Parser :: Unexpected literal")
)

// Инициализация парсера
func newParser(source []byte) *parser {
	return &parser{src: source, line: 1, linePos: 1}
}

// Структура метки для позиции в парсере, которую можно поставить на позицию в парсере для ориентировки или возможного отката
type mark struct {
	pos, line, linePos int
}

// Структура парсера
type parser struct {
	src                []byte  // Срез исходных данных
	pos, line, linePos int     // Информация о позиции в срезе, номере строки, позиции в строке
	_mark              mark    // Метка
	execList           []token // Срез результирующих токенов
}

// Получение среза от установленной метки до текущей позиции со смещением
func (s *parser) markVal(rOffset int) []byte { return s.src[s._mark.pos : s.pos-rOffset] }

// Получение строки от установленной метки со смещением
func (s *parser) markValString(rOffset int) string { return string(s.markVal(rOffset)) }

func (s *parser) parseDocument() (err error) {
	var exec token
	for err == nil && !s.isEndDocument() {
		if exec, err = s.parseToEndLine(); err == nil {
			if exec == nil {
				s.incPos()
			} else if !exec.IsExecutable() {
				err = fmt.Errorf("Executable token expected")
			} else {
				s.execList = append(s.execList, exec)
			}
		} else {
			return
		}
	}
	return
}

// Работа парсера до конца текущей строки
func (s *parser) parseToEndLine() (res token, err error) {
	fmt.Println("=========================================================", s.pos, len(s.src))
	s.passSpaces()
	s.setupMark()
	fmt.Println(!s.isEndLine())
	for !s.isEndLine() && !s.isEndDocument() {
		if opType := checkOpType(s.char()); opType != opUndefined {
			switch opType {
			case opSet:
				{
					if res, err = initSet(s.markVal(0), s); err != nil {
						err = s.setupError(err.Error())
					}
					return
				}
			}
		} else {
			s.incPos()
		}
	}
	s.rollbackMark(0)
	res, err = initVal(s)
	return
}

// Пропуск встречных пробелов
func (s *parser) passSpaces() {
	for !s.isEndDocument() && (s.src[s.pos] == ' ' || s.src[s.pos] == '\t') {
		s.incPos()
	}
}

func (s *parser) passEndLines() {
	for !s.isEndDocument() && (s.src[s.pos] == ' ' || s.src[s.pos] == '\t' || s.src[s.pos] == '\n') {
		s.incPos()
	}
}

func (s *parser) toChar(ch byte) bool {
	for !s.isEndDocument() {
		if s.char() == ch {
			return true
		}
		s.incPos()
	}
	return false
}

// Установка ошибки (к тексту добавляется информация о номере текущей строки и позиции)
func (s *parser) setupError(text string) (err error) {
	err = fmt.Errorf("%s [ line: %d, position: %d ]", text, s.line, s.linePos)
	return
}

// Откат к позиции согласно установленной метки
func (s *parser) rollbackMark(forward int) {
	s.pos, s.line, s.linePos = s._mark.pos, s._mark.line, s._mark.linePos
	for i := 0; i < forward; i++ {
		s.incPos()
	}
}

// Установка метки на текущей позиции
func (s *parser) setupMark() {
	s._mark.pos, s._mark.line, s._mark.linePos = s.pos, s.line, s.linePos
}

// Получение среза от текущей позиции до конца текущей строки (конца документа)
func (s *parser) endLineContent() []byte {
	pos := s.pos
	for pos < len(s.src) && !checkEndLine(s.src[pos]) {
		pos++
	}
	return s.src[s.pos:pos]
}

// Смещение позиции вперед на один символ
func (s *parser) incPos() {
	s.pos++
	if s.pos < len(s.src)-1 {
		if s.src[s.pos] == '\n' {
			s.line, s.linePos = s.line+1, 1
		} else {
			s.linePos++
		}
	}
}

// Получение среза от текущей позиции до конца документа
func (s *parser) availableData() []byte { return s.src[s.pos:] }

func (s *parser) readName() ([]byte, error) {
	s.setupMark()
	if !checkLetter(s.char()) {
		return nil, parseUnexpectedLiteral
	}
	s.incPos()
	for !s.isEndDocument() && (checkLetter(s.char()) || checkNumber(s.char())) {
		s.incPos()
	}
	return s.markVal(0), nil
}

// Получение символа по текущей позиции
func (s *parser) char() byte {
	if s.pos < len(s.src) {
		return s.src[s.pos]
	} else {
		return 0
	}
}

// Проверка соответствия текущей позиции концу документа
func (s *parser) isEndDocument() bool { return len(s.src) <= s.pos }

//func (s *parser) isNextEndDocument() bool { return len(s.src)-1 == s.pos }

// Проверка соответствия текущего символа концу документа или заверщению оператора
func (s *parser) isEndLine() bool {
	if s.isEndDocument() {
		return true
	} else {
		return checkEndLine(s.char())
	}
}

// Проверка соответствия символу завершения оператора (ему соответствует символ конца строки или ";"
func checkEndLine(ch byte) bool {
	return ch == '\n' || ch == ';'
}

//Проверка соответствия первому символу наименования переменной (соответствует регулярке [0-9A-Za-z_])
func checkLetter(ch byte) bool {
	return (ch >= 97 && ch <= 122) || (ch >= 65 && ch <= 90) || ch == 95
}

func checkNumber(ch byte) bool {
	return ch >= 48 && ch <= 57
}

// Шаблон checkLetter, на вхождеие добавлен символ '.'
func checkVarChar(ch byte) bool {
	return checkLetter(ch) || ch == '.'
}
