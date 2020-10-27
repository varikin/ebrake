//+build mage

package main

import (
	"github.com/magefile/mage/sh"
	"strings"
)

func Build() error {
	return run("go build")
}

func run(cmd string) error {
	terms := strings.Fields(cmd)
	if len(terms) == 1 {
		return sh.Run(terms[0])
	} else {
		return sh.Run(terms[0], terms[1:]...)
	}
}
