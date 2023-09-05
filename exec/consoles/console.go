// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoles

type Console interface {
	PollSolicited(messageId int) (response string, hasInput bool)
	PollUnsolicited() (input string, hasInput bool)
	Reset()
	SendReadOnlyMessage(source string, message string)
	SendReadReplyMessage(source string, message string, maxResponseSize int) (messageId int, err error)
	SendStatusMessage(message1 string, message2 string)
}
