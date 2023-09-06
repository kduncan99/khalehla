// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package messages

type ReadOnlyMessage struct {
	source  string   // application sending the message; exec sourced messages are blank
	message []string // message (possibly multi-line)
}

func NewReadOnlyMessage(source string, text string) *ReadOnlyMessage {
	return &ReadOnlyMessage{
		source:  source,
		message: []string{text},
	}
}

func (m *ReadOnlyMessage) GetMessageType() messageType {
	return readOnlyMessageType
}

func (m *ReadOnlyMessage) Serialize() []byte {
	payload := serializeString(m.source)
	payload = append(payload, serializeStringArray(m.message)...)
	return makeMessage(ConsoleMessage, readOnlyMessageType, payload)
}

// -----------------------------------------------------------------------------

type ReadReplyMessage struct {
	messageId      uint32   // unique message identifier for the read-reply message
	source         string   // same as ReadOnlyMessage
	message        []string // same as ReadOnlyMessage
	maxReplyLength uint32   // max characters allowed in the reply
}

func NewReadReplyMessage(messageId uint32, source string, text string, maxReplyLength uint32) *ReadReplyMessage {
	return &ReadReplyMessage{
		messageId:      messageId,
		source:         source,
		message:        []string{text},
		maxReplyLength: maxReplyLength,
	}
}

func (m *ReadReplyMessage) Serialize() []byte {
	payload := serializeUint32(m.messageId)
	payload = append(payload, serializeString(m.source)...)
	payload = append(payload, serializeStringArray(m.message)...)
	payload = append(payload, serializeUint32(m.maxReplyLength)...)
	return makeMessage(ConsoleMessage, readReplyMessageType, payload)
}

// -----------------------------------------------------------------------------

type StatusMessage struct {
	message1 string
	message2 string
}

func NewStatusMessage(text1 string, text2 string) *StatusMessage {
	return &StatusMessage{
		message1: text1,
		message2: text2,
	}
}

func (m *StatusMessage) Serialize() []byte {
	payload := serializeString(m.message1)
	payload = append(payload, serializeString(m.message2)...)
	return makeMessage(ConsoleMessage, statusMessageType, payload)
}
