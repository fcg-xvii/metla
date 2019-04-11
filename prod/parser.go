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
	err       error
}

func (s *parseError) IsNil() bool {
	return s.err == nil
}

func (s *parseError) Error() string {
	return fmt.Sprintf("[%v %v:%v] Parse error - %v", s.tplName, s.line, s.pos, s.err)
}

type parser struct {
	*lineman.CodeLine
	tplName  string
	stack    *containers.Stack
	execList []executer
	store    *storage
	varFlag  bool
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
		s.appendExec(execText{s.MarkVal(offset)})
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
	return s.initParseError(0, 0, fmt.Errorf("Err Parse print"))
}

func (s *parser) parseComment() *parseError {
	s.ForwardPos(2)
	line, pos := s.Line(), s.LinePos()
	if s.ToChar('*') && s.NextChar() == '}' {
		s.ForwardPos(2)
		return nil
	}
	return s.initParseError(line, pos, fmt.Errorf("Unclosed comment tag"))
}

func (s *parser) initParseError(line, pos int, err error) *parseError {
	return &parseError{s.tplName, line, pos, err}
}

func (s *parser) parseCode() *parseError {
	s.ForwardPos(2)
	for !s.isEndCode() {
		fmt.Println("STEP")
		if err := s.initCodeVal(); err != nil {
			return err
		}
	}
	return s.initParseError(0, 0, fmt.Errorf("Err Parse code"))
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
		//return newValString(p)
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
			line, pos := s.Line(), s.LinePos()
			if name, check := s.ReadName(); check {
				if keyword, check := getKeywordConstructor(string(name)); check {
					return keyword(s)
				} else {
					//fmt.Println("NAME", string(s.Char()))
					return newValName(s, line, pos, string(name))
				}
			} else {
				return s.initParseError(s.Line(), s.Pos(), fmt.Errorf("Name read error..."))
			}
		}
	}
	return s.initParseError(s.Line(), s.Pos(), fmt.Errorf("Unexpected synbol %c", s.Char()))
}
