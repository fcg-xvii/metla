package metla

import (
	"fmt"

	"github.com/fcg-xvii/containers"
	"github.com/fcg-xvii/lineman"
)

func initParser(tplName string, src []byte, root *Metla) *parser {
	return &parser{lineman.NewCodeLine(src), root, tplName, new(containers.Stack), nil, new(storage), false, false, false, 0, 0}
}

type parseError struct {
	tplName   string
	line, pos int
	text      string
}

func (s *parseError) Error() string {
	return fmt.Sprintf("[%v %v:%v] Parse error - %v", s.tplName, s.line, s.pos, s.text)
}

type execError struct {
	tplName   string
	line, pos int
	text      string
}

func (s *execError) Error() string {
	return fmt.Sprintf("[%v %v:%v] Exec error - %v", s.tplName, s.line, s.pos, s.text)
}

type parser struct {
	*lineman.CodeLine
	root         *Metla
	tplName      string
	stack        *containers.Stack
	execList     []executer
	store        *storage
	fieldFlag    bool
	varFlag      bool
	rpnFlag      bool
	cycleLayout  byte
	threadLayout byte
}

func (s *parser) PopExecuters() (list []executer, err *parseError) {
	list = make([]executer, s.stack.Len())
	//fmt.Println(s.stack.Len())
	for i, v := range s.stack.PopAll() {
		if ex, check := v.(executer); check {
			list[i] = ex
		} else {
			return nil, v.(coordinator).parseError("Evaluted but not used")
		}
	}
	return
}

func (s *parser) parseDocument() error {
	for !s.IsEndDocument() {
		if err := s.parseText(); err != nil {
			return err
		}
	}
	if s.cycleLayout > 0 {
		return fmt.Errorf("Unclosed for token")
	}
	if s.threadLayout > 0 {
		thread, _, _ := findThread(s)
		return thread.parseError("Unclosed if tag")
	}
	return nil
}

func (s *parser) appendExec(obj executer) {
	s.execList = append(s.execList, obj)
}

func (s *parser) appendText(offset int) {
	src := s.MarkVal(offset)
	if len(src) > 0 {
		s.appendExec(text{position{s.tplName, s.MarkLine(), s.MarkLinePos()}, s.MarkVal(offset)})
	}
}

func (s *parser) parseText() *parseError {
	s.SetupMark()
	check := s.ToChar('{')
	if check {
		switch s.NextChar() {
		case '%':
			s.appendText(0)
			return s.parseCode()
		case '{':
			s.appendText(0)
			return s.parsePrint()
		case '*':
			s.appendText(1)
			return s.parseComment()
		}
	}
	s.IncPos()
	s.appendText(0)
	return nil
}

func (s *parser) IsEndLine() bool {
	if s.CodeLine.IsEndLine() {
		return true
	}
	return s.isEndCode()
}

func (s *parser) isEndCode() bool {
	return s.Char() == '%' && s.NextChar() == '}'
}

func (s *parser) parsePrint() *parseError {
	if err := newPrint(s); err != nil {
		return err
	}
	s.appendExec(s.stack.Pop().(executer))
	return nil
}

func (s *parser) parseComment() *parseError {
	s.ForwardPos(2)
	line, pos := s.Line(), s.LinePos()
	if s.ToChar('*') && s.NextChar() == '}' {
		s.ForwardPos(2)
		return nil
	}
	return s.initParseError(line, pos, "Unclosed comment tag")
}

func (s *parser) initParseError(line, pos int, text string) *parseError {
	return &parseError{s.tplName, line, pos, text}
}

func (s *parser) flushExec() *parseError {
	if exList, err := s.PopExecuters(); err == nil {
		s.execList = append(s.execList, exList...)
	} else {
		return err
	}
	return nil
}

func (s *parser) parseCode() *parseError {
	line, pos := s.Line(), s.LinePos()
	s.ForwardPos(2)
	for !s.IsEndDocument() {
		s.PassSpaces()
		switch {
		case s.Char() == '/' && s.NextChar() == '/':
			s.ToChar('\n')
		case s.Char() == '/' && s.NextChar() == '*':
			line, pos = s.Line(), s.LinePos()
			s.ForwardPos(2)
			for !s.IsEndDocument() {
				if !s.ToChar('*') {
					return s.initParseError(line, pos, "Unclosed comment")
				} else if s.NextChar() == '/' {
					s.ForwardPos(2)
					break
				} else {
					s.IncPos()
				}
			}
		case s.isEndCode():
			if err := s.flushExec(); err != nil {
				return err
			}
			s.ForwardPos(2)
			if s.IsEndLine() {
				s.IncPos()
			}
			return nil
		case s.IsEndLine():
			if err := s.flushExec(); err != nil {
				return err
			}
			s.IncPos()
		default:
			if err := s.initCodeVal(); err != nil {
				return err
			}
		}
	}
	return s.initParseError(line, pos, "Unclosed code tag")
}

func (s *parser) initCodeVal() *parseError {
	s.PassSpaces()
	//fmt.Println("INIT", string(s.Char()))
	switch s.Char() {
	case '+', '-', '*', '/', '(', '!', '>', '<', '%', '&', '|':
		//fmt.Println(s.fieldFlag)
		if s.fieldFlag {
			return newMethod(s)
		} else {
			if s.Char() == ')' {
				if _, check := s.stack.Peek().(iName); check {
					return newFunction(s)
				}
			}
			return newRPN(s)
		}
		/*if p.Char() == '%' && p.NextChar() == '}' {
			return
		}
		return newValArifmetic(p)*/
	case '"', '\'':
		return newValString(s)
	case ',':
		return newValSet(s)
	case '=':
		if s.NextChar() != '=' {
			return newValSet(s)
		} else {
			return newRPN(s)
		}
	case '{':
		return newObject(s)
	case '[':
		if _, check := s.stack.Peek().(iName); check {
			return newIndex(s)
		} else {
			return newArray(s)
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return newValNumber(s)
	case '.':
		return newField(s)
	default:
		if s.IsLetter() != 0 {
			line, pos := s.Line(), s.LinePos()-1
			if name, check := s.ReadName(); check {
				if keyword, check := getKeywordConstructor(string(name)); check {
					return keyword(s)
				} else {
					if !s.fieldFlag {
						if err := newValName(s, line, pos, string(name)); err == nil {
							s.PassSpaces()
							if s.Char() == '(' || s.Char() == '.' || s.Char() == '[' { // Функция
								return s.initCodeVal()
							}
							return nil
						} else {
							return err
						}
					} else {
						s.stack.Push(initStatic(s, string(name), -len(name)))
						return nil
					}
				}
			} else {
				return s.initParseError(s.Line(), s.Pos(), "Name read error...")
			}
		}
	}
	var message string
	switch s.Char() {
	case '\n', ';':
		message = "Unexpected endline"
	default:
		message = fmt.Sprintf("Unexpected symbol %c", s.Char())
	}
	return s.initParseError(s.Line(), s.Pos(), message)
}

func (s *parser) posObject() position {
	return position{s.tplName, s.Line(), s.Pos()}
}

func (s *parser) resetFlags() func() {
	fieldFlag, varFlag, rpnFlag := s.fieldFlag, s.varFlag, s.rpnFlag
	s.fieldFlag, s.varFlag, s.rpnFlag = false, false, false
	return func() {
		s.fieldFlag, s.varFlag, s.rpnFlag = fieldFlag, varFlag, rpnFlag
	}
}
