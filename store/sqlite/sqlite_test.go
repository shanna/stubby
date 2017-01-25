package sqlite_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/shanna/stubby/store/sqlite"
)

var urls = []struct {
	url string
}{
	{"https://google.com"},
	{"https://techspace.co"},
	{"https://shanehanna.org"},
}

func TestStore(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "stubby")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(tmpfile.Name())

	store, err := sqlite.Open(tmpfile.Name())
	if err != nil {
		t.Error(err)
		return
	}
	defer store.Close()

	for _, test := range urls {
		record, err := store.Create(test.url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		record, err = store.Get(record.ID)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
