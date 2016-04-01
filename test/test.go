package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func srcPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func FixturePath() string {
	return filepath.Join(srcPath(), "fixture")
}

func DiffFS(t *testing.T, path string, got, want http.FileSystem) {
	wantFile, err := want.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	gotFile, err := got.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	wantStat, err := wantFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	gotStat, err := gotFile.Stat()
	if err != nil {
		t.Fatal(err)
	}
	diffFileInfo(t, gotStat, wantStat)
	if !gotStat.IsDir() {
		// @TODO test reading with different buffer sizes, etc.
		wantData, err := ioutil.ReadAll(wantFile)
		if err != nil {
			t.Fatal(err)
		}
		gotData, err := ioutil.ReadAll(gotFile)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(gotData, wantData) {
			t.Fatalf("got=%s want=%s", gotData, wantData)
		}
	}
	// @TODO test full readdir semantics
	wantChildren, err := wantFile.Readdir(0)
	if err != nil {
		t.Fatalf("%s: %#v", path, err)
	}
	gotChildren, err := gotFile.Readdir(0)
	if err != nil {
		t.Fatal(err)
	}
	if g, w := len(gotChildren), len(wantChildren); g != w {
		t.Fatalf("got=%d want=%d", g, w)
	}
	for i, wantChild := range wantChildren {
		gotChild := gotChildren[i]
		diffFileInfo(t, gotChild, wantChild)
		DiffFS(t, filepath.Join(path, gotChild.Name()), got, want)
	}
}

func diffFileInfo(t *testing.T, got, want os.FileInfo) {
	if g, w := got.Name(), want.Name(); g != w {
		t.Fatalf("got=%s want=%s", g, w)
	}
	if g, w := got.IsDir(), want.IsDir(); g != w {
		t.Fatalf("got=%t want=%t", g, w)
	}
	if g, w := got.ModTime(), want.ModTime(); !g.Equal(w) {
		t.Fatalf("got=%s want=%s", g, w)
	}
	if g, w := got.Mode(), want.Mode(); g != w {
		t.Fatalf("got=%s want=%s", g, w)
	}
	if g, w := got.Size(), want.Size(); g != w {
		t.Fatalf("got=%d want=%d", g, w)
	}
}
