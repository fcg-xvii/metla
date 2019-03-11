package metla

import (
	"errors"
	"fmt"

	"github.com/fcg-xvii/lineman"
	"github.com/golang-collections/collections/stack"
)

var (
	parseErrEOF            = errors.New("EOF")
	parseUnexpectedLiteral = errors.New("Parser :: Unexpected literal")
)

func newParser(src []byte, tpl *template, root *Metla) *parser {
	return &parser{lineman.NewCodeLine(src), tpl, root, stack.New()}
}

type parser struct {
	*lineman.CodeLine
	tpl       *template
	root      *Metla
	codeStack *stack.Stack
}

func (s *parser) IsEndCode() bool {
	return s.Char() == '%' && s.NextChar() == '}'
}

func (s *parser) flushTextToken() {
	if content := s.MarkVal(0); len(content) > 0 {
		s.tpl.tokenList = append(s.tpl.tokenList, &tokenText{s.infoRecordFromMark(), content})
	}
}

func (s *parser) parseDocument() (err error) {
	passMark := false
	for !s.IsEndDocument() && err == nil {
		if !passMark {
			s.SetupMark()
		}
		bracketOpen := s.ToChar('{')
		if bracketOpen {
			switch s.NextChar() {
			case '{':
				{
					s.flushTextToken()
					s.ForwardPos(2)
					err = s.parsePrint()
				}
			case '%':
				{
					s.flushTextToken()
					s.ForwardPos(2)
					err = s.parseCode()

				}
			default:
				{
					s.IncPos()
					passMark = true
				}
			}
		} else {
			s.flushTextToken()
		}
	}
	return
}

func (s *parser) parsePrint() error {
	s.PassSpaces()
	if token, err := initVal(s); err == nil {
		s.PassSpaces()
		if !s.PosMatchSlice([]byte("}}")) {
			return errors.New("Document parse error :: Unexpected end of print token")
		} else {
			s.tpl.tokenList = append(s.tpl.tokenList, &tokenPrint{s.infoRecordFromMark(), token})
			s.ForwardPos(2)
			return nil
		}
	} else {
		return err
	}
}

func (s *parser) appendExecToken(t token) {
	if t != nil {
		s.tpl.tokenList = append(s.tpl.tokenList, t)
	}
}

func (s *parser) parseCode() (err error) {
	var t token
	for !s.IsEndDocument() && !s.IsEndCode() {
		if t, err = s.parseExecLine(); err == nil {
			s.appendExecToken(t)
		} else {
			return
		}
		s.IncPos()
	}
	if s.IsEndDocument() {
		err = errors.New("Unclosed code tag")
	}
	s.ForwardPos(2)
	return
}

func (s *parser) parseExecLine() (t token, err error) {
	if t, err = s.parseToEndLine(); err == nil && t != nil && !t.IsExecutable() {
		err = fmt.Errorf("Code parse error :: Evaluted but not used [%s]", t)
	}
	return
}

// Парсим до открытия тега кода
func (s *parser) parseToCode() (err error) {
	passMark := false
	for !s.IsEndDocument() && err == nil {
		if !passMark {
			s.SetupMark()
		}
		bracketOpen := s.ToChar('{')
		if bracketOpen {
			switch s.NextChar() {
			case '{':
				{
					s.flushTextToken()
					s.ForwardPos(2)
					err = s.parsePrint()
				}
			case '%':
				{
					return
				}
			default:
				{
					s.IncPos()
					passMark = true
				}
			}
		} else {
			s.flushTextToken()
		}
	}
	return

}

// Парсим код до закрывающего тега (вызывается конструкторами ключевиков)
func (s *parser) parseCodeToCloseTag(tagName string, parent tokenParent) (err error) {
	var t token
	for !s.IsEndDocument() {
		s.PassSpaces()
		if s.IsEndCode() {
			if err = s.parseToCode(); err != nil {
				return
			}
		}

		s.SetupMark()
		if name, check := s.ReadNameSpaces(); check && string(name) == tagName {
			fmt.Println("Name...")
			return nil
		}
		s.RollbackMark(0)

		if t, err = s.parseExecLine(); err != nil {
			return
		} else {
			fmt.Println(t, s.EndLineContent())
			parent.appendChild(t)
			s.IncPos()
		}
	}
	return fmt.Errorf("Unexpected end of document, [%v] expected", string(tagName))
}

// Парсим код до конца строки
func (s *parser) parseToEndLine() (res token, err error) {
	s.PassSpaces()

	// Если после пропуска пробелов конец строки, нет смысла работать дальше
	if s.IsEndLine() {
		return
	}

	// Возможно первое слово - ключевик команды. Если это так - парсим команду
	s.SetupMark()
	if name, check := s.ReadName(); check {
		if keyword, check := getKeywordConstructor(string(name)); check {
			return keyword(s)
		}
	}
	s.RollbackMark(0)
	for !s.IsEndLine() && !s.IsEndDocument() && !s.IsEndCode() {
		s.PassSpaces()
		if opType := checkOpType(s.Char()); opType != opUndefined {
			switch opType {
			case opSet:
				return initSet(s)
			case opArifmetic:
				return initArifmetic(s)
			}
		} else if s.Char() == ',' {
			s.IncPos()
		} else {
			s.SetupMark()
			if res, err = initVal(s); err == nil {
				s.codeStack.Push(res)
			} else {
				return
			}
		}
	}
	return
}

func (s *parser) ParseArgs(stop byte) (err error) {
	var t token
	s.PassSpaces()

	// Если после пропуска пробелов конец строки, нет смысла работать дальше
	if s.IsEndLine() {
		err = fmt.Errorf("Unexpected end line, '%c' expected", stop)
		return
	}
	for !s.IsEndLine() && !s.IsEndDocument() && !s.IsEndCode() {
		if s.Char() == ',' {
			s.IncPos()
		} else if s.Char() == stop {
			return
		} else {
			s.SetupMark()
			if t, err = initVal(s); err == nil {
				s.codeStack.Push(t)
			} else {
				return
			}
		}
	}
	err = s.positionError(fmt.Sprintf("Unexpected symbol '%c', '%c' expected", s.Char(), stop))
	return
}

func (s *parser) positionError(text string) error {
	return fmt.Errorf("Error [%v %v:%v]: %v", s.tpl.objPath, s.MarkLine(), s.MarkLinePos(), text)
}

func (s *parser) infoRecordFromMark() *rawInfoRecord {
	return &rawInfoRecord{tplName: s.tpl.objPath, line: s.MarkLine(), pos: s.MarkLinePos()}
}
