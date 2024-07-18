package main

import (
	"fmt"
	"os"
	"time"

	"github.com/queueue0/qoreutils/internal/terminal"
)

const minColLen = 3

type colInfo struct {
	validLength bool
	lineLength  int
	cols        []int
}

func newColInfo(maxCols int) []colInfo {
	cols := make([]colInfo, maxCols)

	for i := 0; i < maxCols; i++ {
		cols[i].validLength = true
		cols[i].lineLength = (i + 1) * minColLen
		cols[i].cols = make([]int, i+1)

		for j := 0; j < i; j++ {
			cols[i].cols[j] = minColLen
		}
	}

	return cols
}

func (a *arguments) calculateColumns(files []os.DirEntry) (int, []colInfo) {
	lenFiles := len(files)

	// Ignoring err
	width, _, _ := terminal.GetSize(os.Stdout.Fd())
	if width <= 0 {
		width = 80
	}
	maxPossible := width / minColLen

	if width%minColLen != 0 {
		maxPossible += 1
	}

	var maxCols int
	if a.oneColumn {
		maxCols = 1
	} else if maxPossible > 0 && maxPossible < len(files) {
		maxCols = maxPossible
	} else {
		maxCols = lenFiles
	}

	cols := newColInfo(maxCols)

	for f, file := range files {
		nameLen := a.getModdedNameLen(file)

		for i := 0; i < maxCols; i++ {
			if cols[i].validLength {
				idx := f / ((lenFiles + i) / (i + 1))
				realLen := nameLen
				if idx != i {
					realLen += 2
				}

				if cols[i].cols[idx] < realLen {
					cols[i].lineLength += realLen - cols[i].cols[idx]
					cols[i].cols[idx] = realLen
					cols[i].validLength = cols[i].lineLength < width
				}
			}
		}
	}

	var numCols int
	for numCols = maxCols; 1 < numCols; numCols-- {
		if cols[numCols-1].validLength {
			break
		}
	}

	return numCols, cols
}

func formatDate(t time.Time) string {
	currentYear := time.Now().Year()

	formatString := "Jan _2"

	if t.Year() != currentYear {
		formatString = fmt.Sprintf("%s  2006", formatString)
	} else {
		formatString = fmt.Sprintf("%s 15:04", formatString)
	}

	return t.Format(formatString)
}
