package metla

import (
	"fmt"
	"strings"
)

func isPrefixField(prefix string) bool {
	return strings.IndexByte(prefix, '.') != -1
}

func newParser(source []byte) *parser {
	return &parser{src: source}
}

type parserState byte

const (
	scriptKeywordCheck parserState = iota
)

type mark struct {
	pos, line, linePos int
}

type parser struct {
	src                []byte
	pos, line, linePos int
	_mark              mark
	execList           []token
	state              parserState
}

/*func (s *parser) parseDocment() (err error) {
	var prefix string
	for s.pos < len(s.src) {
		if prefix, err = s.parseVarName(); err != nil {
			return
		}
		var op operator
		if constructor, check := getKeywordConstructor(prefix); check {
			if op, err = constructor(prefix, s); err != nil {
				return
			} else {
				s.execList = append(s.execList, op)
			}
		} else {
			if op, err = s.parseOperator(prefix); err != nil {
				return
			} else {
				s.execList = append(s.execList, op)
			}
		}
	}
	return
}*/

func (s *parser) markVal(rOffset int) []byte       { return s.src[s._mark.pos : s.pos-rOffset] }
func (s *parser) markValString(rOffset int) string { return string(s.markVal(rOffset)) }

func (s *parser) parseToEndLine() (res token, err error) {
	fmt.Println("=========================================================")
	s.setupMark()
	for s.pos < len(s.src) && !s.isEndLine() && !s.isEndDocument() {
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
	var length int
	if res, length, err = initVal(s.markVal(0)); err == nil && length < s.pos {
		s.rollbackMark(length)
	}
	return
}

func (s *parser) passSpaces() {
	for !s.isEndDocument() && (s.src[s.pos] == ' ' || s.src[s.pos] == '\t') {
		s.incPos()
	}
}

func (s *parser) setupError(text string) (err error) {
	err = fmt.Errorf("%s [ line: %d, position: %d ]", text, s.line, s.linePos)
	return
}

func (s *parser) rollbackMark(forward int) {
	s.pos, s.line, s.linePos = s._mark.pos, s._mark.line, s._mark.linePos
	for i := 0; i < forward; i++ {
		s.incPos()
	}
}

func (s *parser) setupMark() {
	s._mark.pos, s._mark.line, s._mark.linePos = s.pos, s.line, s.linePos
}

func (s *parser) endLineContent() []byte {
	pos := s.pos
	for pos < len(s.src) && !checkEndLine(s.src[pos]) {
		pos++
	}
	return s.src[s.pos:pos]
}

func (s *parser) incPos() {
	s.pos++
	if s.pos < len(s.src)-1 {
		if s.src[s.pos] == '\n' {
			s.line, s.linePos = s.line+1, 0
		} else {
			s.linePos++
		}
	}
}

func (s *parser) availableData() []byte { return s.src[s.pos:] }
func (s *parser) char() byte            { return s.src[s.pos] }
func (s *parser) isEndDocument() bool   { return len(s.src) == s.pos }
func (s *parser) isEndLine() bool {
	if s.isEndDocument() {
		return true
	} else {
		return checkEndLine(s.char())
	}
}

///////////////////////////////////////////////////////////////////

func checkEndLine(ch byte) bool {
	return ch == '\n' || ch == ';'
}
