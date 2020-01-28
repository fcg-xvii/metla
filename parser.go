package metla

import (
	"bytes"
	"fmt"
)

type tokenType uint8

const (
	tokenTypeUndefined tokenType = iota
	tokenTypeText
	tokenTypePrint
	tokenTypeComment
	tokenTypeExec
	tokenTypeLiteral
)

func (s tokenType) closeTag() byte {
	switch s {
	case tokenTypePrint:
		return '}'
	case tokenTypeExec:
		return '%'
	case tokenTypeComment:
		return '*'
	default:
		return 0
	}
}

func (s tokenType) String() string {
	switch s {
	case tokenTypeUndefined:
		return "Undefined"
	case tokenTypeText:
		return "Text"
	case tokenTypePrint:
		return "Print"
	case tokenTypeComment:
		return "Comment"
	case tokenTypeExec:
		return "Exec"
	case tokenTypeLiteral:
		return "Literal"
	default:
		return ""
	}
}

type tError struct {
	alias string
	line  int
	pos   int
	text  string
}

func (s tError) Error() string {
	return fmt.Sprintf("%v :: Parse error [%v:%v] : %v", s.alias, s.line, s.pos, s.text)
}

func parseBytes(src []byte, alias string) (res []byte, err error) {
	var rb bytes.Buffer
	lp, lastPos, lastLine := lineParser{src: src}, 0, 0

	flushTokenText := func(fromToken bool) {
		if lastPos != lp.Pos() {
			offset := -1
			if fromToken {
				offset = 1
			}
			rb.Write([]byte("flush('"))
			rb.Write(bytes.ReplaceAll(bytes.ReplaceAll(lp.Block(lastPos, offset), []byte("'"), []byte("\\'")), []byte("\n"), []byte("\\n")))
			rb.Write([]byte("');"))
			lastPos, lastLine = lp.Pos(), lp.Line()
		}
	}

	initToken := func() error {
		val, tType := lp.Val(), tokenTypeUndefined
		switch val {
		case '{':
			tType = tokenTypePrint
		case '*':
			tType = tokenTypeComment
		case '%':
			tType = tokenTypeExec
		default:
			return nil
		}
		flushTokenText(true)
		for lp.FindNext(tType.closeTag()) {
			if lp.Next() && lp.Val() == '}' {
				switch tType {
				case tokenTypePrint:
					{
						rb.Write([]byte("flush("))
						rb.Write(lp.Block(lastPos+2, 2))
						rb.Write([]byte(");"))
					}
				case tokenTypeExec:
					rb.Write(lp.Block(lastPos+2, 2))
					rb.Write([]byte(";"))
				}
				lp.Next()
				lastPos, lastLine = lp.Pos(), lp.Line()
				return nil
			}
		}
		return tError{
			alias: alias,
			line:  lastLine,
			pos:   lastPos,
			text:  fmt.Sprintf("Not found close token for '%v}'", string(tType.closeTag())),
		}
	}

	for lp.CheckNext() {
		if lp.FindNext('{') && lp.Next() {
			if err = initToken(); err != nil {
				return
			}
		}
	}
	flushTokenText(false)
	res = rb.Bytes()
	return
}

type lineParser struct {
	src  []byte
	pos  int
	line int
}

func (s *lineParser) CheckNext() bool {
	return s.pos < len(s.src)-1
}

func (s *lineParser) Next() bool {
	if s.CheckNext() {
		s.pos++
		if s.src[s.pos] == '\n' {
			s.line++
		}
		return true
	}
	return false
}

func (s *lineParser) Block(start int, offsetRight int) []byte { return s.src[start : s.pos-offsetRight] }
func (s *lineParser) Line() int                               { return s.line }
func (s *lineParser) Pos() int                                { return s.pos }
func (s *lineParser) Val() byte                               { return s.src[s.pos] }

func (s *lineParser) FindNext(val byte) bool {
	for {
		if s.Val() == val {
			return true
		}
		if !s.Next() {
			return false
		}
	}
}
