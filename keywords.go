package metla

type keywordConstructor func(prefix string, parser *parser) (operator, error)

var (
	keywords = make(map[string]keywordConstructor)
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}
