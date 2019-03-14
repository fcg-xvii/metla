package metla

type keywordConstructor func(*parser) error

var (
	keywords = make(map[string]keywordConstructor)
)

func getKeywordConstructor(name string) (result keywordConstructor, check bool) {
	result, check = keywords[name]
	return
}
