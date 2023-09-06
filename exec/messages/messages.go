// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package messages

type messageCategory uint32

const (
	SystemMessage  messageCategory = 0
	ConsoleMessage messageCategory = 1
)

type messageType uint32

const (
	heartbeatMessageType messageType = 1

	readOnlyMessageType    messageType = 101
	readReplyMessageType   messageType = 102
	resetMessageType       messageType = 103
	statusMessageType      messageType = 104
	solicitedMessageType   messageType = 105
	unsolicitedMessageType messageType = 106
)

const identifier = 2200

func serializeByte(value uint32) []byte {
	return []byte{byte(value)}
}

func serializeUint32(value uint32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(value >> 24)
	bytes[1] = byte((value >> 16) & 0xFF)
	bytes[2] = byte((value >> 8) & 0xFF)
	bytes[3] = byte(value)
	return bytes
}

func serializeString(value string) []byte {
	bytes := serializeUint32(uint32(len(value)))
	bytes = append(bytes, []byte(value)...)
	return bytes
}

func serializeStringArray(array []string) []byte {
	bytes := serializeUint32(uint32(len(array)))
	for _, s := range array {
		bytes = append(bytes, serializeString(s)...)
	}
	return bytes
}

type Message interface {
	GetMessageType() messageType
	Serialize() []byte
}

func makeMessage(msgCategory messageCategory, msgType messageType, payload []byte) []byte {
	msg := serializeUint32(identifier)
	msg = append(msg, serializeUint32(uint32(len(payload)+16))...)
	msg = append(msg, serializeUint32(uint32(msgCategory))...)
	msg = append(msg, serializeUint32(uint32(msgType))...)
	msg = append(msg, payload...)
	return msg
}
