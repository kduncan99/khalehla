// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"khalehla/kexec"
	"os"
	"time"
)

func main() {
	e := kexec.NewExec()
	_ = e.InitialBoot(true)
	time.Sleep(2 * time.Second)
	e.Dump(os.Stdout)
}
