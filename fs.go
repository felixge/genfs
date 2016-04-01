package genfs

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// NewFS returns a new FS for the given files.
func NewFS(files ...*File) *FS {
	fs := &FS{
		l: files,
		m: make(map[string]*File, len(files)),
	}
	for _, file := range files {
		file.fs = fs
		fs.m[file.path] = file
	}
	return fs
}

// Dir returns a new FS for the given path.
func Dir(path string) (*FS, error) {
	files, err := Files(path, FilterNone)
	if err != nil {
		return nil, err
	}
	return NewFS(files...), nil
}

// FS is an in-memory file system that implements http.FileSystem.
type FS struct {
	m map[string]*File
	l []*File
}

// check interface compliance
var _ = http.FileSystem(&FS{})

// Open is part of http.FileSystem.
func (f *FS) Open(name string) (http.File, error) {
	name = path.Clean(name)
	file := f.m[name]
	if file == nil {
		// @TODO Path error?
		return nil, os.ErrNotExist
	}
	clone := &File{}
	*clone = *file
	return clone, nil
}

// NewFile returns a new in-memory file with the given properties.
func NewFile(path, name string, size int64, mode os.FileMode, t time.Time, isDir bool, data []byte) *File {
	return &File{
		path:    path,
		name:    name,
		size:    size,
		mode:    mode,
		modTime: t,
		isDir:   isDir,
		data:    data,
	}
}

// File is an in-memory file that implements http.File and os.FileInfo.
type File struct {
	path    string
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	data    []byte
	fs      *FS
}

// check interface compliance
var _ = os.FileInfo(&File{})
var _ = http.File(&File{})

// Close is part of http.File.
func (f *File) Close() error {
	// @TODO return err when closing twice?
	return nil
}

// Readdir is part of http.File.
func (f *File) Read(p []byte) (int, error) {
	n := copy(p, f.data)
	f.data = f.data[n:]
	if len(f.data) == 0 {
		return n, io.EOF
	} else {
		return n, nil
	}
}

// Seek is part of http.File.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	// @TODO(fg) implement
	return 0, nil
}

// Stat is part of http.File.
func (f *File) Stat() (os.FileInfo, error) {
	// @TODO(fg) implement
	return f, nil
}

// Readdir is part of http.File.
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	// @TODO(fg) support full Readdir semantics, including EOF, etc.
	if count != 0 {
		return nil, errors.New("only count=0 is supported right now")
	}
	var files []os.FileInfo
	// TODO(fg) make this faster using binary search or a tree structure
	for _, file := range f.fs.l {
		if path.Clean(path.Join(f.path, file.name)) == file.path {
			files = append(files, file)
		}
	}
	return files, nil
}

// IsDir is part of os.FileInfo.
func (f *File) IsDir() bool { return f.isDir }

// ModTime is part of os.FileInfo.
func (f *File) ModTime() time.Time { return f.modTime }

// Mode is part of os.FileInfo.
func (f *File) Mode() os.FileMode { return f.mode }

// Name is part of os.FileInfo.
func (f *File) Name() string { return f.name }

// Size is part of os.FileInfo.
func (f *File) Size() int64 { return f.size }

// Sys is part of os.FileInfo.
func (f *File) Sys() interface{} { return nil }

// Path returns the path of the file. It's needed for WriteSource.
func (f *File) Path() string { return f.path }

// String returns the file data as a string. It's needed for WriteSource.
func (f *File) String() string { return string(f.data) }
