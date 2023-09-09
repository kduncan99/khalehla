package main

import (
	"fmt"
	"khalehla/exec/consoles"
	"time"
)

func main() {
	c := consoles.NewNetConsole("127.0.0.1", 2200)
	for true {
		time.Sleep(5 * time.Second)
		msg := fmt.Sprintf("Time: %v", time.Now())
		fmt.Printf("Sending: %v\n", msg)
		c.SendReadOnlyMessage(msg.NewReadOnlyMessage("MAPPER", msg))
	}
}
