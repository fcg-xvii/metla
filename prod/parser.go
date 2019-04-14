package prod

import (
	"fmt"

	"github.com/fcg-xvii/containers"
	"github.com/fcg-xvii/lineman"
)

func initParser(tplName string, src []byte) *parser {
	return &parser{lineman.NewCodeLine(src), tplName, new(containers.Stack), nil, new(storage), false}
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
	tplName  string
	stack    *containers.Stack
	execList []executer
	store    *storage
	varFlag  bool
}

func (s *parser) PopExecuters() (list []executer, err *parseError) {
	list = make([]executer, s.stack.Len())
	fmt.Println(s.stack.Len())
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
	return nil
}

func (s *parser) appendExec(obj executer) {
	s.execList = append(s.execList, obj)
}

func (s *parser) appendText(offset int) {
	src := s.MarkVal(offset)
	if len(src) > 0 {
		s.appendExec(execText{&position{s.tplName, s.MarkLine(), s.MarkLinePos()}, s.MarkVal(offset)})
	}
}

func (s *parser) parseText() *parseError {
	s.SetupMark()
	check := s.ToChar('{')
	if check {
		switch s.NextChar() {
		case '%':
			s.appendText(1)
			return s.parseCode()
		case '{':
			s.appendText(1)
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
	return s.initParseError(0, 0, "Err Parse print")
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

func (s *parser) parseCode() *parseError {
	line, pos := s.Line(), s.LinePos()
	s.ForwardPos(2)
	for !s.IsEndDocument() {
		if err := s.initCodeVal(); err != nil {
			return err
		}
		s.PassSpaces()
		if s.IsEndLine() {
			if exList, err := s.PopExecuters(); err == nil {
				s.execList = append(s.execList, exList...)
				fmt.Println(s.execList)
				if s.isEndCode() {
					s.ForwardPos(2)
					return nil
				} else if s.isEndCode() {
					s.IncPos()
				}
			} else {
				return err
			}
		}
	}
	return s.initParseError(line, pos, "Unclosed code tag")
}

func (s *parser) initCodeVal() *parseError {
	s.PassSpaces()
	fmt.Println("INIT", string(s.Char()))
	if s.IsEndLine() {
		s.IncPos()
		return nil
	}
	switch s.Char() {
	case '+', '-', '*', '/', '(', '!', '>', '<', '%', '&', '|':
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
			return newValArifmetic(s)
		}
	case '{':
		//val, err = newValObject(p)
	case '[':
		//val, err = newValArray(p)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return newValNumber(s)
	case '.':
		/*if !p.fieldFlag {
			p.fieldCommand = true
			val, err = newValField(p)
		}*/
	/*case ':':
	if s.NextChar() == '=' {
		return newValSet(s)
	}*/
	default:
		if s.IsLetter() != 0 {
			line, pos := s.Line(), s.LinePos()-1
			if name, check := s.ReadName(); check {
				if keyword, check := getKeywordConstructor(string(name)); check {
					fmt.Println("KEYWORD!!!!!!!!!!!!", string(name), keyword)
					return keyword(s)
				} else {
					//fmt.Println("NAME", string(s.Char()))
					return newValName(s, line, pos, string(name))
				}
			} else {
				return s.initParseError(s.Line(), s.Pos(), "Name read error...")
			}
		}
	}
	return s.initParseError(s.Line(), s.Pos(), fmt.Sprintf("Unexpected synbol %c", s.Char()))
}
