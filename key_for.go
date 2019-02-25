package metla

import (
	_ "errors"
)

/*func rangeVal(p *parser) (token, error) {
	if !checkValInt(p.EndLineContent()) {
		return nil, errors.New("For (range) parse error :: Unexpected range token type, expected integer")
	}
	return newValInt(p)
}

func keywordFor(p *parser) (res token, err error) {
	var min, max []token
	p.PassSpaces()
	if min, err = rangeVal()
	return nil, nil
}

type keyForCount struct {
	min, max token
}*/
