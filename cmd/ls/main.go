package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

func main() {
	args := os.Args[1:]
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var files []os.DirEntry
	if len(args) == 0 {
		files, err = os.ReadDir(cwd)
	} else {
		files, err = os.ReadDir(args[0])
	}
	if err != nil {
		panic(err)
	}

	tmp := []os.DirEntry{}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
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

	printGrid(files)
}

func printGrid(files []os.DirEntry) {
	cols, colInfo := calculateColumns(files)
	rows := len(files) / cols
	if len(files)%cols != 0 {
		rows += 1
	}
	lineFmt := colInfo[cols-1]

	for row := 0; row < rows; row++ {
		col := 0
		fidx := row

		for {
			nameLen := len([]rune(files[fidx].Name()))
			if col == 0 {
				fmt.Print(" ")
			}
			if strings.Contains(files[fidx].Name(), " ") {
				nameLen += 1
				fmt.Printf("\b'%s'", files[fidx].Name())
			} else {
				fmt.Printf("%s", files[fidx].Name())
			}

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
