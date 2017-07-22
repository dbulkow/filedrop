package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
)

func filehash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open for hash: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Fprintf(os.Stderr, "io copy for hash: %v\n", err)
		os.Exit(1)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	htmlfile := os.Args[len(os.Args)-1]
	if htmlfile == "" {
		fmt.Fprintf(os.Stderr, "no html file specified")
		os.Exit(1)
	}

	parts := strings.Split(htmlfile, ".")
	gofile := parts[0] + "_" + parts[1] + ".go"

	in, err := os.Open(htmlfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open htmlfile: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(gofile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create go file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Fprintf(out, "package main\n\n")
	fmt.Fprintf(out, "const %s_etag = \"%s\"\n\n", parts[0], filehash(htmlfile))
	fmt.Fprintf(out, "const %s = `", parts[0])

	if _, err := io.Copy(out, in); err != nil {
		fmt.Fprintf(os.Stderr, "io copy: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(out, "`\n")
}
