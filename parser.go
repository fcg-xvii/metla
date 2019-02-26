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

func (s *parser) parseDocument() (err error) {
	for !s.IsEndDocument() {
		s.MarkPos()
		s.ToChar('{')
		switch s.NextChar() {
		case '{':
			fmt.Println("PRINT_DOC")
			return
		}
	}
	err = fmt.Errorf("OOOOO")
	return
}

func (s *parser) parsePrint() error {
	for !s.IsEndDocument() {
		s.PassSpaces()
		if token, err := initVal(p); err == nil {
			s.PassSpaces()
			if !s.PosMatchSlice('}}') {
				return errors.New("Document parse error :: Unexpected end of print token")	
			} else {
				
			}
		} else {
			return err
		}
	}
	return nil
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
