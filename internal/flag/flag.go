package flag

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var boolFlags = make(map[string]*bool)
var Args []string

func BoolFlag(name string, ptr *bool) {
	boolFlags[name] = ptr
}

func Parse() error {
	var err error
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flags := strings.Split(strings.Trim(arg, "-"), "")

			for _, flag := range flags {
				if _, ok := boolFlags[flag]; ok {
					*boolFlags[flag] = true
				} else {
					if err == nil {
						err = errors.New(fmt.Sprintf("Unknown flag: %s", flag))
					}
				}
			}
		} else {
			Args = append(Args, arg)
		}
	}

	return err
}
