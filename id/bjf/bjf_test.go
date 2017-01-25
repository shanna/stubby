package bjf_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/shanna/stubby/id/bjf"
)

var ids = []struct {
	id, token string
}{
	{"42", "ecTsv0"},
	{"1", "ecTswn"},
	{"65538", "ecTbtm"},
	{"123456789", "ebqelR"},
	{strconv.Itoa(math.MaxInt32), "c9qJA3"},
}

func TestBjf(t *testing.T) {
	f, err := bjf.New([]byte("the secret sauce"))
	if err != nil {
		t.Errorf("New failed")
		return
	}

	for _, test := range ids {
		token, err := f.Encode(test.id)
		if err != nil {
			t.Errorf("Expected no error on encode but got '%s'", err)
			return
		}
		if token != test.token {
			t.Errorf("Expected token '%s' but got '%s'", test.token, token)
			return
		}

		id, err := f.Decode(token)
		if err != nil {
			t.Errorf("Expected no error on decode but got '%s'", err)
			return
		}
		if id != test.id {
			t.Errorf("Expected bijective rount trip '%s' but got '%s'", test.id, id)
			return
		}
	}
}
