// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"os"

	"khalehla/old/packUtil"
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
		os.Exit(1)
	}
}
