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

func parseFlag(f string) error {
	if _, ok := boolFlags[f]; ok {
		*boolFlags[f] = true
	} else {
		return errors.New(fmt.Sprintf("Unknown flag: %s", f))
	}

	return nil
}

func Parse() error {
	var err error
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "--") {
			flag := strings.Trim(arg, "-")
			e := parseFlag(flag)
			if e != nil && err == nil {
				err = e
			}
		} else if strings.HasPrefix(arg, "-") {
			flags := strings.Split(strings.Trim(arg, "-"), "")

			for _, flag := range flags {
				e := parseFlag(flag)
				if e != nil && err == nil {
					err = e
				}
			}
		} else {
			Args = append(Args, arg)
		}
	}

	return err
}
