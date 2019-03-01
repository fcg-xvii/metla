package metla

import (
	"errors"
	"fmt"

	"github.com/fcg-xvii/lineman"
)

var (
	parseErrEOF            = errors.New("EOF")
	parseUnexpectedLiteral = errors.New("Parser :: Unexpected literal")
)

func newParser(src []byte, tpl *template, root *Metla) *parser {
	return &parser{lineman.NewCodeLine(src), tpl, root}
}

type parser struct {
	*lineman.CodeLine
	tpl  *template
	root *Metla
}

func (s *parser) IsEndCode() bool {
	return s.Char() == '%' && s.NextChar() == '}'
}

func (s *parser) flushTextToken() {
	if content := s.MarkVal(0); len(content) > 0 {
		s.tpl.tokenList = append(s.tpl.tokenList, &tokenText{content})
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
			s.tpl.tokenList = append(s.tpl.tokenList, &tokenPrint{token})
			s.ForwardPos(2)
			return nil
		}
	} else {
		return err
	}
}

func (s *parser) parseCode() (err error) {
	var t token
	for !s.IsEndDocument() && !s.IsEndCode() {
		if t, err = s.parseToEndLine(); err != nil {
			return
		} else {
			fmt.Println("++++++++++++++++++++++++", t)
			s.tpl.tokenList = append(s.tpl.tokenList, t)
			s.IncPos()
		}
	}
	if s.IsEndDocument() {
		err = errors.New("Unclosed code tag")
	}
	return
}

func (s *parser) parseToEndLine() (res token, err error) {
	fmt.Println("=========================================================")
	s.PassSpaces()
	s.SetupMark()
	for !s.IsEndLine() && !s.IsEndDocument() && !s.IsEndCode() {
		if opType := checkOpType(s.Char()); opType != opUndefined {
			switch opType {
			case opSet:
				{
					if res, err = initSet(s.MarkVal(0), s); err != nil {
						err = s.InitError(err.Error())
					}
					return
				}
			}
		} else {
			s.IncPos()
		}
	}
	s.RollbackMark(0)
	fmt.Println("INIT_VAL")
	res, err = initVal(s)
	return
}

func (s *parser) parseCallBack(callback func(*parser) bool) (res token, err error) {
	s.SetupMark()
	for callback(s) {
		s.PassSpaces()
		if opType := checkOpType(s.Char()); opType != opUndefined {
			fmt.Println("OPT", opType)
			return
			/*switch opType {
				case opSet:
			}*/
		} else {
			fmt.Println("S_POS")
			if name, check := s.ReadName(); check {
				if keyword, check := getKeywordConstructor(string(name)); check {
					return keyword(s)
				}
			} else {
				s.IncPos()
			}
		}
	}
	return
}

func (s *parser) parseTempalate() (res token, err error) {
	return
}
