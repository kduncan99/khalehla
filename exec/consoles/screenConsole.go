// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoles

import (
	"fmt"
	"strings"
)

// ScreenConsole implements something closely resembling a UTS60 console.
// Better features include a separate input area to prevent the operator from typing into the display area,
// and a re-callable input history much like one would see in a modern shell.
type ScreenConsole struct {
}

const displayRows = 24
const columns = 80
const inputHistorySize = 20

var inputHistory = make([]string, inputHistorySize)
var inputHistoryCount = 0

type readReplyMessage struct {
	source             string
	message            string
	maxReplyCharacters int
}

var readReplyMessages = map[int]*readReplyMessage{}

// -----------------------------------------------------------------------------

func (c *ScreenConsole) PollSolicited(messageId int) (response string, hasInput bool) {
	return "", false // TODO
}

func (c *ScreenConsole) PollUnsolicited() (input string, hasInput bool) {
	return "", false // TODO
}

func (c *ScreenConsole) Reset() {
	// TODO
}

func (c *ScreenConsole) SendReadOnlyMessage(source string, message string) {
	// TODO
}

func (c *ScreenConsole) SendReadReplyMessage(source string, message string, maxResponseSize int) (messageId int, err error) {
	return 0, fmt.Errorf("foo") // TODO
}

func (c *ScreenConsole) SendStatusMessage(message1 string, message2 string) {
	// TODO
}

// -----------------------------------------------------------------------------

func NewScreenConsole() *ScreenConsole {
	return &ScreenConsole{
		//	TODO
	}
}

// -----------------------------------------------------------------------------

func (c *ScreenConsole) addInputHistoryEntry(input string) {
	upper := strings.ToUpper(input)
	for _, entry := range inputHistory {
		if strings.ToUpper(entry) == upper {
			return
		}
	}

	if inputHistoryCount == inputHistorySize {
		inputHistory = inputHistory[1:]
	}
	inputHistory = append(inputHistory, input)
	inputHistoryCount++
}
