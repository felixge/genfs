package genfs

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/felixge/genfs/test"
)

func TestFS(t *testing.T) {
	cmd := exec.Command("sh", "-c", "find . && echo && pwd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	path := test.FixturePath()
	fmt.Printf("%s\n", path)
	got, err := Dir(path)
	if err != nil {
		t.Fatal(err)
	}
	want := http.Dir(path)
	test.DiffFS(t, ".", got, want)
}
