// khalehla Project
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package channels

// A Channel is a layer between an InputOutputProcessor and a Device.
// It has a goroutine which constantly monitors managed devices, and converts channel programs in main storage
// into IO requests which are sent to the devices, and monitored for completion.
type Channel interface {
	// TODO
}
