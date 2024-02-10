// Khalehla Project
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package consoleMgr

import (
	"fmt"
	"khalehla/kexec/types"
	"khalehla/pkg"
	"log"
	"strings"
	"sync"
	"time"
)

type Console interface {
	ClearReadReplyMessage(messageId int) error
	Close()
	PollSolicitedInput() (*string, int, error)
	PollUnsolicitedInput() (*string, error)
	IsConnected() bool
	Reset() error
	SendReadOnlyMessage(text string) error
	SendSystemMessages(text1 string, text2 string) error
	SendReadReplyMessage(text string, messageId int, maxChars int) error
}

// ConsoleManager handles all things related to console interaction
type ConsoleManager struct {
	consoles        map[pkg.Word36]Console
	terminateThread bool
	threadIsActive  bool
	mutex           sync.Mutex
}

// IsActive indicates whether the goRoutine is active
func (mgr *ConsoleManager) IsActive() bool {
	return mgr.threadIsActive
}

// Reset clears out any artifacts left over by a previous exec session,
// and prepares the console for normal operations
func (mgr *ConsoleManager) Reset() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	if mgr.threadIsActive {
		mgr.terminateThread = true
		for mgr.threadIsActive {
			time.Sleep(1 * time.Second)
		}
	}

	if mgr.consoles == nil {
		// create a single new std console
		mgr.consoles = make(map[pkg.Word36]Console)
		mgr.consoles[0] = NewStandardConsole()
	} else {
		// reset all the existing consoles
		for consId, cons := range mgr.consoles {
			err := cons.Reset()
			if err != nil {
				log.Printf("ConsMgr: Deleting dead console %v", consId.ToStringAsFieldata())
				delete(mgr.consoles, consId)
			}
		}
	}

	go mgr.thread()
}

func (mgr *ConsoleManager) SendReadOnlyMessage(rce *types.RunControlEntry, message string) {
	var text string
	if rce.IsExec {
		text = message
	} else {
		text = fmt.Sprintf("%v*%v", strings.TrimSpace(rce.RunId.ToStringAsFieldata()), text)
	}

	// TODO log it
	// TODO put it in the tailsheet (only if not IsExec)

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for consId, cons := range mgr.consoles {
		err := cons.SendReadOnlyMessage(text)
		if err != nil {
			log.Printf("ConsMgr: Deleting dead console %v", consId.ToStringAsFieldata())
			delete(mgr.consoles, consId)
		}
	}
}

// Stop is invoked when the exec is stopping... for any reason. It tells the goRoutine to stop.
func (mgr *ConsoleManager) Stop() {
	mgr.terminateThread = true
}

// thread is the main routine for the console manager goRoutine
func (mgr *ConsoleManager) thread() {
	mgr.threadIsActive = true

	// TODO

	mgr.threadIsActive = false
}
