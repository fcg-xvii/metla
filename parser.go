package metla

import (
	"errors"
	"fmt"

	"github.com/fcg-xvii/lineman"
)

type tokenContainer interface {
	pushToken(token)
}

type parseMethod func(p *parser) error

var (
	parseErrEOF            = errors.New("EOF")
	parseUnexpectedLiteral = errors.New("Parser :: Unexpected literal")
)

func newParser(src []byte, tpl *template, root *Metla) *parser {
	return &parser{lineman.NewCodeLine(src), tpl, root, newCodeStack(), true}
}

type parser struct {
	*lineman.CodeLine
	tpl       *template
	root      *Metla
	stack     *codeStack
	textState bool
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
	s.stack.Push(initSplitter())
}

func (s *parser) flushStack() {
	s.tpl.tokenList = append(s.tpl.tokenList, s.stack.Flush()...)
}
