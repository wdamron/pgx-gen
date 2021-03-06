package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wdamron/astx"
	"github.com/wdamron/pgx-gen"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		Usage()
		os.Exit(1)
	}
	path := os.Args[1]
	if path == "" {
		Usage()
		os.Exit(1)
	}
	af, err := astx.ParseFile(path)
	if err != nil {
		Err(err)
		os.Exit(1)
	}
	var outpath string
	if len(os.Args) == 3 {
		outpath = os.Args[2]
	} else {
		outpath = strings.TrimSuffix(path, filepath.Ext(path)) + "_pgxgen.go"
	}

	f := pgxgen.NewFile(af)
	gen, err := f.Gen()
	if err != nil {
		Err(err)
		os.Exit(1)
	}

	out, err := os.Create(outpath)
	if err != nil {
		Err(err)
		os.Exit(1)
	}
	defer out.Close()
	_, err = out.Write(gen)
	if err != nil {
		Err(err)
		os.Exit(1)
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tpqx-gen filepath [outpath]\n\n")
	fmt.Fprintf(os.Stderr, "Defaults:\n")
	fmt.Fprintf(os.Stderr, "\toutpath: filepath + \"_pgxgen.go\"\n\n")
}

func Err(err error) {
	fmt.Fprintln(os.Stderr, err)
}
