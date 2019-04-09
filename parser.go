package metla

import (
	"errors"
	"fmt"

	"github.com/fcg-xvii/lineman"
	"github.com/golang-collections/collections/stack"
)

type parseMethod func(p *parser) error

var (
	parseErrEOF            = errors.New("EOF")
	parseUnexpectedLiteral = errors.New("Parser :: Unexpected literal")
)

func newParser(src []byte, tpl *Template, root *Metla) *parser {
	return &parser{lineman.NewCodeLine(src), tpl, root, stack.New(), stack.New(), true, false, false}
}

type parser struct {
	*lineman.CodeLine
	tpl          *Template
	root         *Metla
	stack        *stack.Stack
	openStack    *stack.Stack
	textState    bool
	fieldFlag    bool
	fieldCommand bool
}

func (s *parser) parseDocument() (err error) {
	for !s.IsEndDocument() {
		var method parseMethod
		switch s.Char() {
		case '{':
			{
				switch s.NextChar() {
				case '{':
					method = newValPrint
				case '*':
					method = newValComment
				case '%':
					method = newValCode
				}
			}
		default:
			method = newValText
		}
		if err = method(s); err != nil {
			return
		}
	}
	s.flushStack()
	if s.openStack.Len() > 0 {
		openInfo := s.openStack.Pop().(openFlag)
		err = openInfo.info.fatalError(fmt.Sprintf("Unclosed tag %v", openInfo.tagName))
	}
	return
}

func (s *parser) markError(text string) error {
	return fmt.Errorf("Error [%v %v:%v]: %v", s.tpl.objPath, s.MarkLine(), s.MarkLinePos(), text)
}

func (s *parser) positionError(text string) error {
	return fmt.Errorf("Error [%v %v:%v]: %v", s.tpl.objPath, s.Line(), s.LinePos(), text)
}

func (s *parser) infoRecordFromMark() *rawInfoRecord {
	return &rawInfoRecord{tplName: s.tpl.objPath, line: s.MarkLine(), pos: s.MarkLinePos()}
}

func (s *parser) infoRecordFromPos() *rawInfoRecord {
	return &rawInfoRecord{tplName: s.tpl.objPath, line: s.Line(), pos: s.LinePos()}
}

func (s *parser) pushSplitter() {
	//s.stack.Push(initSplitter())
}

func (s *parser) stackToSlice() []interface{} {
	res := make([]interface{}, 0, s.stack.Len())
	for i := s.stack.Len() - 1; i >= 0; i-- {
		res = append(res, s.stack.Pop())
	}
	return res
}

func (s *parser) flushStack() {
	s.tpl.tokenList = append(s.tpl.tokenList, s.stackToSlice()...)
}

func (s *parser) readStackVal() []interface{} {
	if _, check := s.stack.Peek().(*execMarker); check {
		childsCounter := 0
		var res []interface{}
		for s.stack.Len() > 0 {
			val := s.stack.Pop()
			switch val.(type) {
			case *execMarker:
				res = append(res, val)
				childsCounter++
			case *execCommand:
				res = append(res, val)
				childsCounter--
				if childsCounter == 0 {
					return res
				}
			default:
				res = append(res, val)
			}
		}
		return res
	} else {
		return []interface{}{s.stack.Pop()}
	}
}
