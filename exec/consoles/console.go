// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoles

import "khalehla/exec/messages"

type Console interface {
	ClearReadReplyMessage(messageId int) (err error)
	IsReady() bool
	PollSolicitedInput(messageId int) (response string, hasInput bool, err error)
	PollUnsolicitedInput() (input string, hasInput bool)
	Reset()
	SendReadOnlyMessage(message *messages.ReadOnlyMessage)
	SendReadReplyMessage(message *messages.ReadReplyMessage)
	SendStatusMessage(message *messages.StatusMessage)
}
