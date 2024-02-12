// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoleMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type readReplyTracker struct {
	trackerId    int
	message      *types.ConsoleReadReplyMessage
	replyConsole *types.ConsoleIdentifier // this is the console which is currently supposed to answer
	hasReply     bool
	isCanceled   bool
	messageId    int // from the Console
	retryLater   bool
}

// ConsoleManager handles all things related to console interaction
type ConsoleManager struct {
	exec             types.IExec
	consoles         map[types.ConsoleIdentifier]types.Console
	primaryConsole   types.Console
	primaryConsoleId types.ConsoleIdentifier
	terminateThread  bool
	threadStarted    bool
	threadStopped    bool
	mutex            sync.Mutex
	queuedReadOnly   []*types.ConsoleReadOnlyMessage
	queuedReadReply  map[int]*readReplyTracker
}

func NewConsoleManager(exec types.IExec) *ConsoleManager {
	return &ConsoleManager{
		exec: exec,
	}
}

// IsDone indicates whether the goRoutine is active
func (mgr *ConsoleManager) IsDone() bool {
	return mgr.threadStopped
}

// CloseManager is invoked when the exec is stopping... for any reason. It tells the goRoutine to threadStop.
func (mgr *ConsoleManager) CloseManager() {
	mgr.threadStop()
}

func (mgr *ConsoleManager) InitializeManager() {
	mgr.consoles = make(map[types.ConsoleIdentifier]types.Console)
	mgr.queuedReadOnly = make([]*types.ConsoleReadOnlyMessage, 0)
	mgr.queuedReadReply = make(map[int]*readReplyTracker)

	mgr.primaryConsole = NewStandardConsole()
	mgr.primaryConsoleId = types.ConsoleIdentifier(pkg.NewFromStringToFieldata("SYSCON", 1)[0])
	mgr.consoles[mgr.primaryConsoleId] = mgr.primaryConsole

	mgr.threadStart()
}

// ResetManager clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *ConsoleManager) ResetManager() {
	mgr.threadStop()

	mgr.mutex.Lock()
	if mgr.consoles == nil {
		// create a single new std console
		mgr.consoles = make(map[types.ConsoleIdentifier]types.Console)
		mgr.consoles[0] = NewStandardConsole()
	} else {
		// reset all the existing consoles
		for consId, cons := range mgr.consoles {
			err := cons.Reset()
			if err != nil {
				mgr.dropConsole(consId)
			}
		}
	}
	mgr.mutex.Unlock()

	mgr.threadStart()
}

// SendReadOnlyMessage queues a RO message and returns immediately.
// The ConsoleManager thread will handle actually sending the message to all the consoles if/as appropriate.
func (mgr *ConsoleManager) SendReadOnlyMessage(message *types.ConsoleReadOnlyMessage) {
	// Log it and put it in the RCE tail sheet (unless it is the Exec)
	log.Printf("%v*%v", message.Source.RunId.ToStringAsFieldata(), message.Text)
	if !message.Source.IsExec {
		message.Source.PrintToTailSheet(message.Text)
	}

	mgr.mutex.Lock()
	mgr.queuedReadOnly = append(mgr.queuedReadOnly, message)
	mgr.mutex.Unlock()
}

// SendReadReplyMessage queues a read-reply message and waits for the response.
// The wait is terminated if the RCE goes into contingency mode, or if the exec stops.
// During the waiting period, the ConsoleManager thread will send the message, then poll for a reply as necessary.
func (mgr *ConsoleManager) SendReadReplyMessage(message *types.ConsoleReadReplyMessage) error {
	// Log it and put it in the RCE tail sheet (unless it is the Exec)
	log.Printf("*-%v:%v", message.Source.RunId.ToStringAsFieldata(), message.Text)
	if !message.Source.IsExec {
		message.Source.PrintToTailSheet(message.Text)
	}

	tracker := mgr.newReadReplyTracker(message)
	for !tracker.hasReply && !tracker.isCanceled {
		time.Sleep(25 * time.Millisecond)
	}

	mgr.mutex.Lock()
	delete(mgr.queuedReadReply, tracker.trackerId)
	mgr.mutex.Unlock()

	if tracker.isCanceled {
		return fmt.Errorf("read reply message was canceled")
	}

	return nil
}

// SendSystemMessages does a best-effort job at sending the system messages to all the consoles.
// Whether the consoles do anything with them is up to them.
// We do not log any of these.
func (mgr *ConsoleManager) SendSystemMessages(message1 string, message2 string) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for consId, cons := range mgr.consoles {
		err := cons.SendSystemMessages(message1, message2)
		if err != nil {
			mgr.dropConsole(consId)
		}
	}
}

// checkForReadOnlyMessages pops the next message from the read only queue (if any)
// and sends it on its way.
func (mgr *ConsoleManager) checkForReadOnlyMessages() bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if len(mgr.queuedReadOnly) > 0 {
		msg := mgr.queuedReadOnly[0]
		mgr.queuedReadOnly = mgr.queuedReadOnly[1:]

		// Construct output text
		text := ""
		if !msg.DoNotEmitRunId {
			text = msg.Source.RunId.ToStringAsFieldata() + "*"
		}
		text += msg.Text

		// If it has routing, try to send it to the indicated console
		if msg.Routing != nil {
			cons, ok := mgr.consoles[*msg.Routing]
			if ok {
				err := cons.SendReadOnlyMessage(text)
				if err == nil {
					return true
				}

				// lose the console, and drop through
				mgr.dropConsole(*msg.Routing)
			}

			// routing was desired, but could not be fulfilled.
			// send it to the primary console instead
			_ = mgr.primaryConsole.SendReadOnlyMessage(text)
			return true
		}

		// No routing - send it to all the consoles
		for consId, cons := range mgr.consoles {
			err := cons.SendReadOnlyMessage(text)
			if err != nil {
				mgr.dropConsole(consId)
			}
		}

		return true
	}

	return false
}

// checkForReadReplyMessages pops the next message *which does not have a reply console*
// from the read reply queue and posts it as follows:
//
//	If there is routing, send it to that console as read-reply
//	If that fails, send it to the primary console as read-reply
//	If there is NO routing, send it to the primary console as read-reply, and copy it to other consoles as read-only
func (mgr *ConsoleManager) checkForReadReplyMessages() bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for _, tracker := range mgr.queuedReadReply {
		if tracker.replyConsole == nil && !tracker.retryLater {
			// Construct output text
			text := ""
			if !tracker.message.DoNotEmitRunId {
				text = tracker.message.Source.RunId.ToStringAsFieldata() + "*"
			}
			text += tracker.message.Text

			// If it has routing, try to send it to the indicated console
			var console types.Console = nil
			var consoleId types.ConsoleIdentifier
			if tracker.message.Routing != nil {
				cons, ok := mgr.consoles[*tracker.message.Routing]
				if ok {
					console = cons
					consoleId = *tracker.message.Routing
				}
			}

			if console == nil {
				console = mgr.primaryConsole
				consoleId = mgr.primaryConsoleId
			}

			msgId, err := console.SendReadReplyMessage(text, tracker.message.MaxReplyLength)
			if err != nil {
				// probably a message id overflow, although a dead console cannot be discounted
				// set it up for a retry - if the console is really dead, it will soon be dropped.
				// At that point we'll go back through this code but we'll take a different path.
				tracker.retryLater = true
				continue // see if there's a different message we can send
			} else {
				tracker.replyConsole = &consoleId
				tracker.messageId = msgId
			}

			// Now *maybe* send copies to all the other consoles
			if tracker.message.Routing == nil {
				for consId, cons := range mgr.consoles {
					if consId != consoleId {
						err := cons.SendReadOnlyMessage(text)
						if err != nil {
							mgr.dropConsole(consId)
						}
					}
				}
			}

			return true
		}
	}

	return false
}

func (mgr *ConsoleManager) checkForSolicitedInput() bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for consId, cons := range mgr.consoles {
		reply, msgId, err := cons.PollSolicitedInput()
		if err != nil {
			// kill to console, and reschedule any RR messages waiting on that console
			mgr.dropConsole(consId)
		} else if reply != nil {
			// we have a reply - find the tracker for this message
			trackerId, err := mgr.findTrackerFor(consId, msgId)
			if err != nil {
				log.Printf("ConsMgr:Received reply for a message we are not tracking")
				continue
			}

			tracker := mgr.queuedReadReply[trackerId]
			if tracker.hasReply {
				log.Printf("ConsMgr:Received reply for a message we already have a reply for")
				continue
			}
			if len(*reply) > tracker.message.MaxReplyLength {
				_ = cons.SendReadOnlyMessage("REPLY TOO LONG - RE-ENTER")
				continue
			}

			if !tracker.message.DoNotLogReply {
				log.Printf("ConsMgr:%v %v", msgId, reply)
				tracker.message.Source.PrintToTailSheet(fmt.Sprintf("%v %v", msgId, reply))
			}

			tracker.message.Reply = *reply
			tracker.hasReply = true
			delete(mgr.queuedReadReply, trackerId)
			return true
		}
	}

	return false
}

func (mgr *ConsoleManager) checkForUnsolicitedInput() bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	for consId, cons := range mgr.consoles {
		input, err := cons.PollUnsolicitedInput()
		if err != nil {
			mgr.dropConsole(consId)
		} else if input != nil {
			// send the raw input to the exec, and let it deal with parsing issues
			mgr.exec.HandleKeyIn(consId, *input)
			return true
		}
	}

	return false
}

// dropConsole is invoked whenever any higher level code gets an error response from a Console.
// call under lock
func (mgr *ConsoleManager) dropConsole(consoleId types.ConsoleIdentifier) {
	consId := pkg.Word36(consoleId)
	log.Printf("ConsMgr: Deleting unreponsive console %v", consId.ToStringAsFieldata())
	delete(mgr.consoles, consoleId)

	for _, tracker := range mgr.queuedReadReply {
		if tracker.replyConsole != nil && *tracker.replyConsole == consoleId {
			tracker.replyConsole = nil
		}
	}
}

// call under lock
func (mgr *ConsoleManager) findTrackerFor(consoleId types.ConsoleIdentifier, messageId int) (int, error) {
	for trackerId, tracker := range mgr.queuedReadReply {
		if tracker.replyConsole != nil && *tracker.replyConsole == consoleId && tracker.messageId == messageId {
			return trackerId, nil
		}
	}

	return 0, fmt.Errorf("tracker not found")
}

// newReadReplyTracker generates a readReplyTracker for the given message and queues it up
func (mgr *ConsoleManager) newReadReplyTracker(message *types.ConsoleReadReplyMessage) *readReplyTracker {
	tracker := &readReplyTracker{}
	tracker.message = message
	tracker.hasReply = false
	tracker.isCanceled = false
	tracker.replyConsole = nil

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	tracker.trackerId = 1
	_, ok := mgr.queuedReadReply[tracker.trackerId]
	for ok {
		tracker.trackerId++
		_, ok = mgr.queuedReadReply[tracker.trackerId]
	}

	mgr.queuedReadReply[tracker.trackerId] = tracker
	return tracker
}

// thread is the main routine for the console manager goRoutine
func (mgr *ConsoleManager) thread() {
	mgr.threadStarted = true

	retryCounter := 0
	for !mgr.terminateThread {
		result := mgr.checkForReadOnlyMessages()
		result = result || mgr.checkForReadReplyMessages()
		result = result || mgr.checkForSolicitedInput()
		result = result || mgr.checkForUnsolicitedInput()

		if !result {
			time.Sleep(250 * time.Millisecond)
			retryCounter++
			if retryCounter >= 8 {
				mgr.mutex.Lock()
				for _, tracker := range mgr.queuedReadReply {
					tracker.retryLater = false
				}
				mgr.mutex.Unlock()
			}
		}
	}

	// cancel all outstanding read-reply messages
	mgr.mutex.Lock()
	for trackerId, tracker := range mgr.queuedReadReply {
		if !tracker.isCanceled && !tracker.hasReply {
			tracker.isCanceled = true
			delete(mgr.queuedReadReply, trackerId)
		}
	}
	mgr.mutex.Unlock()

	mgr.threadStarted = false
}

func (mgr *ConsoleManager) threadStart() {
	mgr.terminateThread = false
	if !mgr.threadStarted {
		go mgr.thread()
		for !mgr.threadStarted {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *ConsoleManager) threadStop() {
	if mgr.threadStarted {
		mgr.terminateThread = true
		for !mgr.threadStopped {
			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (mgr *ConsoleManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vConsoleManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, mgr.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, mgr.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, mgr.terminateThread)
	_, _ = fmt.Fprintf(dest, "%v  Consoles:\n", indent)
	primaryStr := ""
	for consId, cons := range mgr.consoles {
		if consId == mgr.primaryConsoleId {
			primaryStr = " (Primary)"
		} else {
			primaryStr = ""
		}

		consoleId := pkg.Word36(consId)
		_, _ = fmt.Fprintf(dest, "%v    %v%v\n", indent, consoleId.ToStringAsFieldata(), primaryStr)
		cons.Dump(dest, indent+"  ")
	}

	_, _ = fmt.Fprintf(dest, "%v  QueuedReadOnly:\n", indent)
	for _, msg := range mgr.queuedReadOnly {
		str := "[" + msg.Source.RunId.ToStringAsFieldata() + "] " + msg.Text

		if msg.DoNotEmitRunId {
			str += " !emitRunId "
		}

		if msg.RunId != nil {
			str += " RunId:" + *msg.RunId + " "
		}

		if msg.Routing != nil {
			consId := pkg.Word36(*msg.Routing)
			str += " Routing:" + consId.ToStringAsFieldata() + " "
		}
		_, _ = fmt.Fprintf(dest, "%v    %v\n", indent, str)
	}

	_, _ = fmt.Fprintf(dest, "%v  QueuedReadReply:\n", indent)
	for tid, tracker := range mgr.queuedReadReply {
		str := fmt.Sprintf("tid:%v hasReply:%v canceled:%v retry:%v",
			tid, tracker.hasReply, tracker.isCanceled, tracker.retryLater)
		if tracker.replyConsole != nil {
			consId := pkg.Word36(*tracker.replyConsole)
			str += " repCons:" + consId.ToStringAsFieldata()
		}
		_, _ = fmt.Fprintf(dest, "%v    %v\n", indent, str)

		msg := tracker.message
		str = "[ " + msg.Source.RunId.ToStringAsFieldata() + "] " + msg.Text

		if msg.DoNotEmitRunId {
			str += " !emitRunId "
		}

		if msg.RunId != nil {
			str += " RunId:" + *msg.RunId + " "
		}

		if msg.Routing != nil {
			consId := pkg.Word36(*msg.Routing)
			str += " Routing:" + consId.ToStringAsFieldata() + " "
		}

		str += fmt.Sprintf("max:%v   ", msg.MaxReplyLength)

		if msg.DoNotLogReply {
			str += " !logReply"
		} else {
			str += " reply:" + msg.Reply
		}

		_, _ = fmt.Fprintf(dest, "%v      %v\n", indent, str)
	}
}
