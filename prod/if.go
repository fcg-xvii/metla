package prod

import (
	"fmt"
)

func init() {
	keywords["if"] = newIf
}

func newIf(p *parser) *parseError {

}
