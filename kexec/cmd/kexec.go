// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"khalehla/kexec"
	"khalehla/kexec/config"
	"time"
)

func main() {
	cfg := &config.Configuration{}
	e := kexec.NewExec(cfg)
	_ = e.InitialBoot(true)
	for !e.GetStopFlag() {
		time.Sleep(250 * time.Millisecond)
	}
}
