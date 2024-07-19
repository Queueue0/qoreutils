package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/queueue0/qoreutils/internal/flag"
)

type arguments struct {
	oneColumn     bool
	longList      bool
	showAll       bool
	showAlmostAll bool
	hasQuotes     bool
}

func main() {
	a := arguments{}
	flag.BoolFlag("a", &a.showAll)
	flag.BoolFlag("A", &a.showAlmostAll)
	flag.BoolFlag("all", &a.showAll)
	flag.BoolFlag("almost-all", &a.showAlmostAll)
	flag.BoolFlag("1", &a.oneColumn)
	flag.BoolFlag("l", &a.longList)
	err := flag.Parse()
	if err != nil {
		// TODO: print usage instead
		panic(err)
	}
	args := flag.Args

	if a.showAlmostAll && !a.showAll {
		a.showAll = true
	}

	var files []os.DirEntry
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

	if a.showAll && !a.showAlmostAll {
		cwdfi, err := os.Lstat(".")
		if err == nil {
			cwdde := fs.FileInfoToDirEntry(cwdfi)
			files = append(files, cwdde)
		}

		parentfi, err := os.Lstat("..")
		if err == nil {
			files = append(files, fs.FileInfoToDirEntry(parentfi))
		}
	}

	tmp := []os.DirEntry{}
	for _, file := range files {
		if a.showAll || !strings.HasPrefix(file.Name(), ".") {
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

	if a.longList {
		a.printLongList(files)
	} else {
		a.printGrid(files)
	}
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

type llFile struct {
	mode      string
	hardlinks string
	username  string
	groupname string
	size      string
	modtime   string
	name      string
}

func (a *arguments) printLongList(files []os.DirEntry) {
	llFiles := []llFile{}
	maxLink := 0
	maxSize := 0
	for _, f := range files {
		info, err := f.Info()
		if err != nil {
			panic(err)
		}
		links := uint64(0)
		size := uint64(0)
		var username string
		var groupname string
		if sys := info.Sys(); sys != nil {
			if stat, ok := sys.(*syscall.Stat_t); ok {
				links = uint64(stat.Nlink)
				size = uint64(stat.Size)
				uid := strconv.Itoa(int(stat.Uid))
				gid := strconv.Itoa(int(stat.Gid))
				u, err := user.LookupId(uid)
				if err != nil {
					panic(err)
				}
				g, err := user.LookupGroupId(gid)
				if err != nil {
					panic(err)
				}
				username = u.Username
				groupname = g.Name
			}
		}
		modtime := formatDate(info.ModTime())

		llf := llFile{info.Mode().String(), strconv.Itoa(int(links)), username, groupname, strconv.Itoa(int(size)), modtime, a.getModdedName(f)}
		llFiles = append(llFiles, llf)

		linklen := len([]rune(llf.hardlinks))
		sizelen := len([]rune(llf.size))
		if linklen > maxLink {
			maxLink = linklen
		}

		if sizelen > maxSize {
			maxSize = sizelen
		}
	}

	for _, llf := range llFiles {
		linklen := len([]rune(llf.hardlinks))
		sizelen := len([]rune(llf.size))

		linkpad := ""
		for linklen < maxLink {
			linkpad += " "
			linklen++
		}

		sizepad := ""
		for sizelen < maxSize {
			sizepad += " "
			sizelen++
		}

		fmt.Printf("%s %s%s %s %s %s%s %s %s\n", llf.mode, linkpad, llf.hardlinks, llf.username, llf.groupname, sizepad, llf.size, llf.modtime, llf.name)
	}
}
