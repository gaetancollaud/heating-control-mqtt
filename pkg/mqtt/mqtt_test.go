package mqtt

import (
	"testing"
)

func expect(t *testing.T, result string, expect string, msg string) {
	if expect != result {
		t.Errorf("%s Expected='%s' but got '%s'", msg, expect, result)
	}
}
