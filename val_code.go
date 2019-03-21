package metla

import _ "fmt"

func newValCode(p *parser) error {
	p.ForwardPos(2)
	for !p.IsEndDocument() {
		p.PassSpaces()
		switch p.Char() {
		case ';', '\n':
			p.flushStack()
			p.IncPos()
		case '%':
			if p.NextChar() == '}' {
				p.flushStack()
				p.ForwardPos(2)
				return nil
			}
		default:
			if _, err := initCodeVal(p); err != nil {
				return err
			}
		}
	}
	return p.positionError("CODE_ERR")
}

func newValIndex(p *parser) (interface{}, error) {
	return nil, p.positionError("index_error")
}

func newValField(p *parser) (interface{}, error) {
	return nil, p.positionError("field_error")
}
