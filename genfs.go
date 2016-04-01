package genfs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type Filter func(string) bool

var FilterNone = func(string) bool { return false }

// FilterRegexp returns a filter func for the given posix regexp expr, or an
// error.
func FilterRegexp(expr string) (Filter, error) {
	r, err := regexp.CompilePOSIX(expr)
	if err != nil {
		return nil, err
	}
	return func(name string) bool {
		return r.MatchString(name)
	}, nil
}

func Dir(path string) (*FS, error) {
	files, err := Files(path, FilterNone)
	if err != nil {
		return nil, err
	}
	return NewFS(files...), nil
}

func Files(path string, ignore Filter) ([]*File, error) {
	return files(path, path, ignore)
}

func files(path, root string, ignore Filter) ([]*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return nil, err
	}
	if ignore(relPath) {
		return nil, nil
	}
	result := &File{
		path:    relPath,
		name:    stat.Name(),
		isDir:   stat.IsDir(),
		size:    stat.Size(),
		mode:    stat.Mode(),
		modTime: stat.ModTime(),
	}
	if !stat.IsDir() {
		result.data, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return []*File{result}, nil
	}
	children, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}
	results := make([]*File, 0, len(children)+1)
	results = append(results, result)
	for _, child := range children {
		childResults, err := files(filepath.Join(path, child.Name()), root, ignore)
		if err != nil {
			return nil, err
		}
		results = append(results, childResults...)
	}
	return results, nil
}

func WriteSource(w io.Writer, pkg, varName string, files []*File) error {

	return tmpl.Execute(w, struct {
		Package string
		Var     string
		Files   []*File
	}{pkg, varName, files})
}

var tmpl = template.Must(
	template.New("").
		Funcs(map[string]interface{}{
		"timeFormat": func(t time.Time) string {
			return t.Format(time.RFC3339Nano)
		},
	}).
		Parse(strings.TrimSpace(`
package {{.Package}}

import (
	"time"
	"github.com/felixge/genfs"
)

func init () {
	{{.Var}} = genfs.NewFS({{range .Files}}
		genfs.NewFile(
			{{printf "%#v" .Path}},
			{{printf "%#v" .Name}},
			{{printf "%#v" .Size}},
			{{printf "%#v" .Mode}},
			mustTime({{printf "%#v" (timeFormat .ModTime)}}),
			{{printf "%#v" .IsDir}},
			[]byte({{printf "%#v" .String}}),
		),{{end}}
	)
}

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
`)))
