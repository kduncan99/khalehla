// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoleMgr

import (
	"fmt"
	"io"
	"khalehla/kexec"
	"khalehla/pkg"
	"log"
	"sync"
	"time"
)

type readReplyTracker struct {
	trackerId    int
	message      *kexec.ConsoleReadReplyMessage
	replyConsole *kexec.ConsoleIdentifier // this is the console which is currently supposed to answer
	hasReply     bool
	isCanceled   bool
	messageId    int // from the IConsole
	retryLater   bool
}

// ConsoleManager handles all things related to console interaction
type ConsoleManager struct {
	exec             kexec.IExec
	mutex            sync.Mutex
	threadDone       bool
	threadStop       bool
	consoles         map[kexec.ConsoleIdentifier]kexec.IConsole
	primaryConsole   kexec.IConsole
	primaryConsoleId kexec.ConsoleIdentifier
	queuedReadOnly   []*kexec.ConsoleReadOnlyMessage
	queuedReadReply  map[int]*readReplyTracker
}

func NewConsoleManager(exec kexec.IExec) *ConsoleManager {
	return &ConsoleManager{
		exec: exec,
	}
}

// Boot is invoked when the exec is booting
func (mgr *ConsoleManager) Boot() error {
	log.Printf("ConsMgr:Boot")
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	//	TODO shut down all known net consoles

	// reset the consoles list to include only the existing system console
	mgr.consoles = map[kexec.ConsoleIdentifier]kexec.IConsole{
		mgr.primaryConsoleId: mgr.primaryConsole,
	}
	_ = mgr.primaryConsole.Reset()

	// clear the console queues
	mgr.queuedReadOnly = make([]*kexec.ConsoleReadOnlyMessage, 0)
	mgr.queuedReadReply = make(map[int]*readReplyTracker)
	return nil
}

// Close is invoked when the application is terminating
func (mgr *ConsoleManager) Close() {
	log.Printf("ConsMgr:Close")
	mgr.threadStop = true
	for !mgr.threadDone {
		time.Sleep(25 * time.Millisecond)
	}
}

// Initialize is invoked when the application is starting
func (mgr *ConsoleManager) Initialize() error {
	log.Printf("ConsMgr:Initialize")
	mgr.consoles = make(map[kexec.ConsoleIdentifier]kexec.IConsole)
	mgr.queuedReadOnly = make([]*kexec.ConsoleReadOnlyMessage, 0)
	mgr.queuedReadReply = make(map[int]*readReplyTracker)

	mgr.primaryConsole = NewStandardConsole()
	mgr.primaryConsoleId = kexec.ConsoleIdentifier(pkg.NewFromStringToFieldata("SYSCON", 1)[0])
	mgr.consoles[mgr.primaryConsoleId] = mgr.primaryConsole

	// TODO Load net console configuration

	// The console manager thread is always running, although the net facility may not be accepting connections
	go mgr.thread()
	return nil
}

// Stop is invoked when the exec is stopping
func (mgr *ConsoleManager) Stop() {
	log.Printf("ConsMgr:Stop")
	// TODO
	//  We should probably do something here, but I'm not sure what.
}

// SendReadOnlyMessage queues a RO message and returns immediately.
// The ConsoleManager thread will handle actually sending the message to all the consoles if/as appropriate.
func (mgr *ConsoleManager) SendReadOnlyMessage(message *kexec.ConsoleReadOnlyMessage) {
	// Log it and put it in the RCE tail sheet (unless it is the Exec)
	log.Printf("ConsMgr:Queueing %v*%v", message.Source.GetRunId(), message.Text)
	if !message.Source.IsExec() {
		message.Source.PostToTailSheet(message.Text)
	}

	mgr.mutex.Lock()
	mgr.queuedReadOnly = append(mgr.queuedReadOnly, message)
	mgr.mutex.Unlock()
}

// SendReadReplyMessage queues a read-reply message and waits for the response.
// The wait is terminated if the RCE goes into contingency mode, or if the exec stops.
// During the waiting period, the ConsoleManager thread will send the message, then poll for a reply as necessary.
func (mgr *ConsoleManager) SendReadReplyMessage(message *kexec.ConsoleReadReplyMessage) error {
	// Log it and put it in the RCE tail sheet (unless it is the Exec)
	log.Printf("ConsMgr:Queueing n-%v:%v", message.Source.GetRunId(), message.Text)
	if !message.Source.IsExec() {
		message.Source.PostToTailSheet(message.Text)
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
			text = msg.Source.GetRunId() + "*"
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
		if mgr.exec.GetStopFlag() {
			if !tracker.hasReply {
				tracker.isCanceled = true
			}
			continue
		}

		if tracker.replyConsole == nil && !tracker.retryLater {
			// Construct output text
			text := ""
			if !tracker.message.DoNotEmitRunId {
				text = tracker.message.Source.GetRunId() + "*"
			}
			text += tracker.message.Text

			// If it has routing, try to send it to the indicated console
			var console kexec.IConsole = nil
			var consoleId kexec.ConsoleIdentifier
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
				// At that point we'll go back through this code, but we'll take a different path.
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
				consw36 := pkg.Word36(consId)
				log.Printf("ConsMgr:%v %v %v", consw36.ToStringAsFieldata(), msgId, *reply)
				tracker.message.Source.PostToTailSheet(fmt.Sprintf("%v %v", msgId, reply))
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
			km := mgr.exec.GetKeyinManager()
			km.PostKeyin(consId, *input)
			return true
		}
	}

	return false
}

// dropConsole is invoked whenever any higher level code gets an error response from a IConsole.
// call under lock
func (mgr *ConsoleManager) dropConsole(consoleId kexec.ConsoleIdentifier) {
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
func (mgr *ConsoleManager) findTrackerFor(consoleId kexec.ConsoleIdentifier, messageId int) (int, error) {
	for trackerId, tracker := range mgr.queuedReadReply {
		if tracker.replyConsole != nil && *tracker.replyConsole == consoleId && tracker.messageId == messageId {
			return trackerId, nil
		}
	}

	return 0, fmt.Errorf("tracker not found")
}

// newReadReplyTracker generates a readReplyTracker for the given message and queues it up
func (mgr *ConsoleManager) newReadReplyTracker(message *kexec.ConsoleReadReplyMessage) *readReplyTracker {
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
// It runs once across all exec sessions, so it is started during Initialize() and terminated by Close().
func (mgr *ConsoleManager) thread() {
	mgr.threadDone = false

	retryCounter := 0 // we only check retries every 8 times through the loop
	for !mgr.threadStop {
		// TODO check for any new net console connections...
		// 	or should we have a separate coroutine for this?

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
	mgr.threadDone = true
}

func (mgr *ConsoleManager) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vConsoleManager ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadDone:   %v\n", indent, mgr.threadDone)
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
		str := "[" + msg.Source.GetRunId() + "] " + msg.Text

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
		str = "[ " + msg.Source.GetRunId() + "] " + msg.Text

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
