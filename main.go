package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func dirTreeRec(out io.Writer, path string, printFiles bool, level int, mainPrefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var filteredFiles []fs.DirEntry
	for _, entry := range files {
		if entry.IsDir() || printFiles {
			filteredFiles = append(filteredFiles, entry)
		}
	}

	var prefix strings.Builder
	var prefixedPath strings.Builder

	for i, entry := range filteredFiles {
		prefix.Reset()
		prefix.WriteString(mainPrefix)
		prefixedPath.Reset()

		if i == len(filteredFiles)-1 {
			prefix.WriteString("└───")
		} else {
			prefix.WriteString("├───")
		}

		prefixedPath.WriteString(prefix.String())
		prefixedPath.WriteString(entry.Name())

		if entry.IsDir() {
			fmt.Fprintln(out, prefixedPath.String())
			newPrefix := mainPrefix
			if i == len(filteredFiles)-1 {
				newPrefix += "\t"
			} else {
				newPrefix += "│\t"
			}
			dirTreeRec(out, filepath.Join(path, entry.Name()), printFiles, level+1, newPrefix)

		} else if printFiles {
			info, err := entry.Info()
			if err != nil {
				return fmt.Errorf("error geting file info: %v", err)
			}
			fileSize := info.Size()

			if fileSize == 0 {
				prefixedPath.WriteString(" (empty)")
			} else {
				prefixedPath.WriteString(fmt.Sprintf(" (%db)", fileSize))
			}
			fmt.Fprintln(out, prefixedPath.String())
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeRec(out, path, printFiles, 0, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
