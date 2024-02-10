package main

import (
	"khalehla/kexec"
)

func main() {
	e := kexec.Exec{}
	_ = e.InitialBoot(true)
}
