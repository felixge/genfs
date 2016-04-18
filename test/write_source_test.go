package test

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/felixge/genfs"
)

var genFS http.FileSystem

func TestWriteSource(t *testing.T) {
	files, err := genfs.Files(FixturePath())
	if err != nil {
		t.Fatal(err)
	}
	want := genfs.NewFS(files...)
	if genFS == nil {
		path := filepath.Join(srcPath(), "genfs_fs_test.go")
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		defer os.Remove(path)
		if err := genfs.WriteSource(file, "test", "genFS", files); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("go", "test", ".")
		cmd.Dir = srcPath()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	} else {
		DiffFS(t, ".", genFS, want)
	}
}
