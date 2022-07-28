package randstring

import (
	"fmt"
	"testing"
)

func TestRandStringRunes(t *testing.T) {
	rs := RandString(10)
	fmt.Printf("random string: %s\n", rs)
}
