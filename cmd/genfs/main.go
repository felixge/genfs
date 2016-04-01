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
		ignoreF = flag.String("ignore", "/\\.", "Regular expression for files and directory paths to ignore.")
		pkgF    = flag.String("pkg", "main", "Pkg name for the generated file.")
		varF    = flag.String("var", "fs", "Name of the variable to hold the fs.")
	)
	flag.Parse()
	ignore, err := genfs.FilterRegexp(*ignoreF)
	if err != nil {
		fmt.Printf("-ignore: %s\n", err)
	}
	files, err := genfs.Files(flag.Arg(0), ignore)
	if err != nil {
		return err
	}
	return genfs.WriteSource(os.Stdout, *pkgF, *varF, files)
}
