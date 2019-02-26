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

func (s *parser) flushTextToken() {
	if content := s.MarkVal(0); len(content) > 0 {
		s.tpl.tokenList = append(s.tpl.tokenList, &tokenText{content})
	}
}

func (s *parser) parseDocument() (err error) {
	for !s.IsEndDocument() && err == nil {
		s.MarkPos()
		bracketOpen := s.ToChar('{')
		s.flushTextToken()
		if bracketOpen {
			switch s.NextChar() {
			case '{':
				{
					s.ForwardPos(2)
					err = s.parsePrint()
				}
			}
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

/*func (s *parser) parseDocument() (err error) {
	var exec token
	for err == nil && !s.IsEndDocument() {
		if exec, err = s.parseToEndLine(); err == nil {
			if exec == nil {
				s.IncPos()
			} else if !exec.IsExecutable() {
				err = fmt.Errorf("Document parse error :: Executable token expected")
			} else {
				//s.execList = append(s.execList, exec)
				s.tpl.tokenList = append(s.tpl.tokenList, exec)
			}
		} else {
			return
		}
	}
	return
}*/

func (s *parser) parseToEndLine() (res token, err error) {
	fmt.Println("=========================================================")
	s.PassSpaces()
	s.SetupMark()
	for !s.IsEndLine() && !s.IsEndDocument() {
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
