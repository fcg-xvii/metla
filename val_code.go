package metla

func newValCode(p *parser) error {
	return p.positionError("CODE_ERR")
}

func newValSet(p *parser) (interface{}, error) {
	return nil, p.positionError("set_error")
}

func newValIndex(p *parser) (interface{}, error) {
	return nil, p.positionError("index_error")
}

func newValField(p *parser) (interface{}, error) {
	return nil, p.positionError("field_error")
}
