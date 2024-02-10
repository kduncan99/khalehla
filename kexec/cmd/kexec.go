// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package main

import (
	"khalehla/kexec"
)

func main() {
	e := kexec.Exec{}
	_ = e.InitialBoot(true)
}
