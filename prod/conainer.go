package prod

func newField(p *parser) *parseError {
	pos := position{p.tplName, p.Line(), p.LinePos()}
	if p.stack.Len() == 0 {
		return p.initParseError(pos.line, pos.pos, "Expected field owner")
	}

}
