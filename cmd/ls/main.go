package main

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/queueue0/qoreutils/internal/flag"
)

type arguments struct {
	showHidden bool
	hasQuotes  bool
}

func main() {
	a := arguments{}
	flag.BoolFlag("a", &a.showHidden)
	_ = flag.Parse()
	args := flag.Args

	var files []os.DirEntry
	var err error
	if len(args) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		files, err = os.ReadDir(cwd)
		if err != nil {
			panic(err)
		}
	} else {
		files, err = os.ReadDir(args[0])
		if err != nil {
			panic(err)
		}
		if err = os.Chdir(args[0]); err != nil {
			panic(err)
		}
	}

	cwdfi, err := os.Lstat(".")
	if err == nil {
		cwdde := fs.FileInfoToDirEntry(cwdfi)
		files = append(files, cwdde)
	}

	parentfi, err := os.Lstat("..")
	if err == nil {
		files = append(files, fs.FileInfoToDirEntry(parentfi))
	}

	tmp := []os.DirEntry{}
	for _, file := range files {
		if a.showHidden || !strings.HasPrefix(file.Name(), ".") {
			tmp = append(tmp, file)
		}
	}
	files = tmp

	cmp := func(a, b os.DirEntry) int {
		aName, bName := a.Name(), b.Name()
		aName, bName = strings.ToLower(aName), strings.ToLower(bName)

		r, _ := regexp.Compile("[^a-z0-9]+")
		aName, bName = r.ReplaceAllString(aName, ""), r.ReplaceAllString(bName, "")

		if aName < bName {
			return -1
		}

		if aName > bName {
			return 1
		}

		return 0
	}

	slices.SortFunc(files, cmp)

	a.printGrid(files)
}

func (a *arguments) printGrid(files []os.DirEntry) {
	cols, colInfo := a.calculateColumns(files)
	rows := len(files) / cols
	if len(files)%cols != 0 {
		rows += 1
	}
	lineFmt := colInfo[cols-1]

	for row := 0; row < rows; row++ {
		col := 0
		fidx := row

		for {
			nameLen := a.getModdedNameLen(files[fidx])
			if a.hasQuotes && col == 0 {
				fmt.Print(" ")
			}
			fmt.Print(a.getModdedName(files[fidx]))

			fidx += rows
			if fidx >= len(files) {
				break
			}

			maxNameLen := lineFmt.cols[col]
			col += 1

			from, to := nameLen, maxNameLen
			for from < to {
				fmt.Print(" ")
				from++
			}
		}
		fmt.Print("\n")
	}
}
