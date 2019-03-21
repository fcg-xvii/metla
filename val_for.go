package metla

func newValFor(p *parser) (res interface{}, err error) {
	err = p.positionError("ENDFOR")
	return
}
