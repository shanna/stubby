package stubby_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/shanna/stubby"
	"github.com/shanna/stubby/id/bjf"
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

	id, err := bjf.New([]byte("secret"))
	if err != nil {
		t.Error(err)
		return
	}

	shortener, err := stubby.New(store, id)
	if err != nil {
		t.Error(err)
		return
	}

	for _, test := range urls {
		record, err := shortener.Create(test.url, nil)
		if err != nil {
			t.Error(err)
			return
		}

		record, err = shortener.Get(record.ID)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
