package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/felixge/genfs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		includeF = flag.String("include", "", "Regular expression for files and directory paths to ignore.")
		excludeF = flag.String("exclude", "", "Regular expression for files and directory paths to exclude.")
		pkgF     = flag.String("pkg", "main", "Pkg name for the generated file.")
		varF     = flag.String("var", "fs", "Name of the variable to hold the fs.")
	)
	flag.Parse()
	var filters []genfs.Filter
	if *includeF != "" {
		f, err := genfs.IncludeRegexp(*includeF)
		if err != nil {
			return fmt.Errorf("-include: %s\n", err)
		}
		filters = append(filters, f)
	}
	if *excludeF != "" {
		f, err := genfs.IncludeRegexp(*excludeF)
		if err != nil {
			return fmt.Errorf("-exclude: %s\n", err)
		}
		filters = append(filters, f)
	}
	files, err := genfs.Files(flag.Arg(0), filters...)
	if err != nil {
		return err
	}
	return genfs.WriteSource(os.Stdout, *pkgF, *varF, files)
}
