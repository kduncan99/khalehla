// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoleMgr

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type stdConsReadReplyInfo struct {
	messageId         int
	originalMessage   string
	maxResponseLength int
	response          *string
}

// StandardConsole implements a basic (if possibly a little flaky) console using stdin and stdout
type StandardConsole struct {
	mutex                   sync.Mutex
	activeReadReplyMessages map[int]*stdConsReadReplyInfo
	pendingReadReply        *int
	pendingUnsolicitedInput *string
	reader                  *bufio.Reader
	termFlag                bool
}

func NewStandardConsole() *StandardConsole {
	con := &StandardConsole{
		activeReadReplyMessages: make(map[int]*stdConsReadReplyInfo),
		reader:                  bufio.NewReader(os.Stdin),
		termFlag:                false,
	}

	go con.routine()
	return con
}

func (con *StandardConsole) ClearReadReplyMessage(messageId int) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	info, ok := con.activeReadReplyMessages[messageId]
	if !ok {
		return fmt.Errorf("message not found")
	}

	fmt.Printf("~ %v\n", info.originalMessage)
	delete(con.activeReadReplyMessages, messageId)
	return nil
}

func (con *StandardConsole) Close() {
	// shut down - kills the routine loop
	con.termFlag = true
}

func (con *StandardConsole) PollSolicitedInput() (*string, int, error) {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	if con.pendingReadReply != nil {
		rr := con.activeReadReplyMessages[*con.pendingReadReply]
		delete(con.activeReadReplyMessages, rr.messageId)
		con.pendingReadReply = nil
		fmt.Printf("  %v\n", rr.originalMessage)
		return rr.response, rr.messageId, nil
	}

	return nil, 0, nil
}

func (con *StandardConsole) PollUnsolicitedInput() (*string, error) {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	result := con.pendingUnsolicitedInput
	con.pendingUnsolicitedInput = nil
	return result, nil
}

func (con *StandardConsole) IsConnected() bool {
	// we are always connected
	return true
}

func (con *StandardConsole) Reset() error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	con.activeReadReplyMessages = make(map[int]*stdConsReadReplyInfo)
	con.pendingReadReply = nil
	con.pendingUnsolicitedInput = nil

	fmt.Println()
	fmt.Println()
	fmt.Println("*** CONSOLE RESET ***")
	return nil
}

func (con *StandardConsole) SendReadOnlyMessage(text string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	fmt.Printf(" %v\n", text)
	return nil
}

func (con *StandardConsole) SendSystemMessages(string, string) error {
	// we don't display system messages
	return nil
}

// SendReadReplyMessage attempts to send the given message after assigning it a unique identifier
// If there are no available identifiers, we return an error.
// The caller should, in this case, try again later.
func (con *StandardConsole) SendReadReplyMessage(text string, maxChars int) (int, error) {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	msgId := 0
	_, ok := con.activeReadReplyMessages[msgId]
	for ok {
		msgId++
		if msgId > 9 {
			return 0, fmt.Errorf("message id overflow")
		}
		_, ok = con.activeReadReplyMessages[msgId]
	}

	rr := &stdConsReadReplyInfo{
		messageId:         msgId,
		originalMessage:   text,
		maxResponseLength: maxChars,
		response:          nil,
	}

	fmt.Printf("%v-%v\n", msgId, text)
	con.activeReadReplyMessages[msgId] = rr
	return msgId, nil
}

func parseResponse(input string) (bool, int, string, error) {
	if len(input) >= 0 && input[0] >= '0' && input[1] <= '9' {
		msgId := int(input[0] - '0')
		if len(input) == 1 {
			return true, msgId, "", nil
		}
		if input[1] == ' ' {
			return true, msgId, strings.TrimSpace(input[1:]), nil
		}
		return false, 0, "", fmt.Errorf("invalid input")
	}

	return false, 0, "", nil
}

func (con *StandardConsole) pollInput() {
	_ = os.Stdin.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	bytes, _, err := con.reader.ReadLine()
	if !con.termFlag && err == nil {
		con.mutex.Lock()
		defer con.mutex.Unlock()

		// Is there anything pending? If so, ignore this
		if con.pendingUnsolicitedInput != nil || con.pendingReadReply != nil {
			fmt.Println("WAIT - INPUT IGNORED")
			return
		}

		input := string(bytes)
		if len(input) > 0 {
			isReply, msgId, text, err := parseResponse(input)
			if err != nil {
				fmt.Println("INVALID INPUT")
				return
			}

			if isReply {
				// (potential) response to read-reply message
				rr, ok := con.activeReadReplyMessages[msgId]
				if !ok {
					fmt.Println("MESSAGE DOES NOT EXIST")
					return
				}

				if len(text) > rr.maxResponseLength {
					fmt.Println("RESPONSE TOO LONG")
					return
				}

				rr.response = &text
				con.pendingReadReply = &msgId
			} else {
				// unsolicited input
				input = strings.TrimSpace(input)
				con.pendingUnsolicitedInput = &input
			}
		}
	}
}

func (con *StandardConsole) routine() {
	for !con.termFlag {
		con.pollInput()
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("** CONSOLE CLOSED **")
}
