// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package hardware

import "fmt"

func DumpByteArray(buffer []byte) {
	for bx := 0; bx < len(buffer); bx += 32 {
		for by := 0; by < 32; by++ {
			fmt.Printf("%02X ", buffer[bx+by])
			if bx+by == len(buffer) {
				break
			}
		}
		fmt.Println()
	}
}
