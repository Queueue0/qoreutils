package main

import (
	"fmt"
	"os"
	"strings"
)

func (a *arguments) getModdedName(f os.DirEntry) string {
	name := f.Name()

	if strings.Contains(name, " ") {
		if !a.hasQuotes {
			a.hasQuotes = true
		}
		name = fmt.Sprintf("\b'%s'", name)
	}

	return name
}

func (a *arguments) getModdedNameLen(f os.DirEntry) int {
	name := f.Name()

	if strings.Contains(name, " ") {
		if !a.hasQuotes {
			a.hasQuotes = true
		}
		return len([]rune(name)) + 1
	}

	return len([]rune(name))
}
