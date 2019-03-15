package metla

func newValCode(p *parser) error {
	return p.positionError("CODE_ERR")
}

func newValSet(p *parser) error {
	return p.positionError("set_error")
}

func newValIndex(p *parser) error {
	return p.positionError("index_error")
}

func newValField(p *parser) error {
	return p.positionError("field_error")
}
