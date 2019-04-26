package prod

func newFunction(p *parser) *parseError {
	return p.initParseError(10, 10, "Function error test")
}

/*type function struct {

}*/
