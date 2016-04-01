package genfs

import (
	"net/http"
	"testing"

	"github.com/felixge/genfs/test"
)

func TestFS(t *testing.T) {
	path := test.FixturePath()
	got, err := Dir(path)
	if err != nil {
		t.Fatal(err)
	}
	want := http.Dir(path)
	test.DiffFS(t, ".", got, want)
}
