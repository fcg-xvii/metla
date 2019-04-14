package prod

type keywordConstructor func(*parser) *parseError

func init() {
	//keywords["echo"] = keywordEcho
	//keywords["echoln"] = keywordEcholn
	//keywords["print"] = keywordPrint
	//keywords["println"] = keywordPrintln
}

var (
	keywords = map[string]keywordConstructor{
		"nil": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, nil, 3))
			return nil
		}, "true": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, true, 4))
			return nil
		}, "false": func(p *parser) *parseError {
			p.stack.Push(initStatic(p, false, 5))
			return nil
		}, "var": func(p *parser) *parseError {
			if p.varFlag {
				return p.initParseError(p.Line(), p.LinePos()-3, "Unexpected var keyword")
			}
			p.varFlag = true
			return nil
		},
	}
	functions = map[string]interface{}{
		//"len":     coreLen,
		//"defined": coreDefined,
	}
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}
