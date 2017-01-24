package bjf_test

import (
	"testing"

	"github.com/shanna/stubby/id/bjf"
)

// TODO: Table driven. Check a bunch of valid IDs.
var ids = []struct {
	id, token string
}{
	{"42", "ecTsv0"},
	{"1", "ecTswn"},
	{"65538", "ecTbtm"},
	{"123456789", "ebqelR"},
}

func TestBjf(t *testing.T) {
	f, err := bjf.New([]byte("the secret sauce"))
	if err != nil {
		t.Errorf("New failed")
		return
	}

	for _, id := range ids {
		token := f.Encode(id.id)
		if token != id.token {
			t.Errorf("Expected token '%s' but got '%s'", id.token, token)
			return
		}

		if f.Decode(token) != id.id {
			t.Errorf("Encode -> Decode bijective function failure.")
			return
		}
	}
}
