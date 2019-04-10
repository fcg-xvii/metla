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
			p.stack.Push(nil)
			return nil
		}, "true": func(p *parser) *parseError {
			p.stack.Push(true)
			return nil
		}, "false": func(p *parser) *parseError {
			p.stack.Push(false)
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
