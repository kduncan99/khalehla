// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"khalehla/kexec"
	"khalehla/kexec/config"
	"os"
	"time"
)

func main() {
	cfg := config.NewConfiguration()
	cfg.LogIOs = true

	e := kexec.NewExec(cfg)
	e.SetJumpKey(3, true)
	e.SetJumpKey(4, true)
	e.SetJumpKey(7, true)
	e.SetJumpKey(13, true)

	_ = e.InitialBoot(true)
	for !e.GetStopFlag() {
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(2 * time.Second)

	e.Dump(os.Stdout)
}
