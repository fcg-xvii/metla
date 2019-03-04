package metla

import (
	"bytes"
	"fmt"
)

func checkValFor(src []byte) bool {
	return bytes.Index(src, []byte("for")) == 0
}

func newValFor(p *parser) (res token, err error) {
	err = fmt.Errorf("For Test Error...")
	return
}
