package main

import (
	"fmt"
	"khalehla/kexec/packUtil"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		packUtil.DoUsage()
		os.Exit(1)
	}

	var err error
	if args[0] == "prep" {
		err = packUtil.DoPrep(args[1:])
	} else if args[0] == "show" {
		err = packUtil.DoShow(args[1:])
	} else {
		packUtil.DoUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
