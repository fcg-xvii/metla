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

type parser struct {
	src                 []byte
	line, linePos       int
	pos, mark, lineMark int
	execList            []token
	state               parserState
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

func (s *parser) markVal(rOffset int) []byte       { return s.src[s.mark : s.pos-rOffset] }
func (s *parser) markValString(rOffset int) string { return string(s.markVal(rOffset)) }

func (s *parser) parseToEndLine() (res token, err error) {
	fmt.Println("=========================================================")
	s.setupMark()
	for s.pos < len(s.src) && !s.isEndLine() && !s.isEndDocument() {
		if opType := checkOpType(s.char()); opType != opUndefined {
			switch opType {
			case opSet:
				{
					if res, err = initSet(s.markValString(0), s); err != nil {
						err = s.setupError(err.Error())
						return
					}
				}
			}
		} else {
			s.incPos()
		}
	}
	return initVal(s.markVal(0))
}

func (s *parser) passSpaces() {
	for !s.isEndDocument() && (s.src[s.pos] == ' ' || s.src[s.pos] == '\t') {
		s.incPos()
	}
}

/*func (s *parser) parseVarName() (result string, err error) {
	s.passSpaces()
	ch := s.src[s.pos]
	s.setupMark()
	if !checkLetter(ch) {
		err = s.setupError(fmt.Sprintf("Unexpected keyword or variable first (%c)", ch))
		return
	}
	s.incPos()
	ch = s.src[s.pos]
	for s.pos < len(s.src) && (checkVarChar(ch) || checkVarChar(ch)) {
		s.incPos()
		ch = s.src[s.pos]
	}
	result = string(s.src[s.mark:s.pos])
	return
}*/

/*func (s *parser) parseOperator(prefix string) (op operator, err error) {
	s.passSpaces()
	ch := s.src[s.pos]
	opType := checkOpType(ch)
	switch opType {
	case opSet:
		{
			if op, err = initSet(prefix, s); err != nil {
				err = s.setupError(err.Error())
			}
		}
	default:
		{
			return nil, s.setupError(fmt.Sprintf("Unexpected operator (%c)", ch))
		}

	}
	return
}*/

func (s *parser) setupError(text string) (err error) {
	err = fmt.Errorf("%s [ line: %d, position: %d ]", text, s.line, s.linePos)
	return
}

func (s *parser) setupMark() {
	s.mark = s.pos
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
func (s *parser) isEndLine() bool       { ch := s.char(); return ch == '\n' || ch == ';' }
