// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import (
	"fmt"
	"io"
	"khalehla/kexec/nodeMgr"
	"khalehla/kexec/types"
	"strings"
)

/*
FS KEYIN - component DOES NOT EXIST, INPUT IGNORED
FS KEYIN - eqp-mnemonic EQUIPMENT MNEMONIC ILLEGAL, INPUT IGNORED
  (Exec) An equipment mnemonic cannot be entered on any keyin other than the DN,PACK keyin.
FS KEYIN - inhibits INHIBITS ILLEGAL, INPUT IGNORED
  (Exec) Inhibits cannot be entered on any keyin other than the MD keyin.
FS KEYIN NOT ALLOWED - DIRECTORY ID MUST BE dir-id
  (Exec) The directory-id of a pack and a directory-id specified on the keyin are opposite
  (for example, a local pack-id and shared were specified on the keyin).
FS KEYIN - NO UNIT EXISTS IN THE PREMOUNT ONLY STATUS, INPUT IGNORED
FS KEYIN - option OPTION DOES NOT EXIST, INPUT IGNORED
  (Exec) An illegal FS keyin was entered. The variable option is the option that you specified on your FS keyin.
FS NOT ALLOWED UNTIL MASS STORAGE INITIALIZED OR RECOVERED
  (Exec) An FS,PACK keyin is not allowed until the recovery files have been created or restored.
FS,PACK KEYIN ERROR - MHFS IS NOT AVAILABLE
  (Exec) The FS,PACK/SHARED keyin is not allowed because Multi-Host File Sharing (MHFS) is down or not available.
FS,PACK NOT ALLOWED - dir-id IS ILLEGAL DIRECTORY ID

Variations we accept:
	FS,[ CM | DISK | FDISK | MS | PACK | RDISK | TAPE ]
	FS node_name[,...]
	FS,ALL channel_name
*/

type FSKeyinHandler struct {
	exec            types.IExec
	source          types.ConsoleIdentifier
	options         string
	arguments       string
	terminateThread bool
	threadStarted   bool
	threadStopped   bool
}

func NewFSKeyinHandler(exec types.IExec, source types.ConsoleIdentifier, options string, arguments string) *FSKeyinHandler {
	return &FSKeyinHandler{
		exec:            exec,
		source:          source,
		options:         strings.ToUpper(options),
		arguments:       strings.ToUpper(arguments),
		terminateThread: false,
		threadStarted:   false,
		threadStopped:   false,
	}
}

func (kh *FSKeyinHandler) Abort() {
	kh.terminateThread = true
}

func (kh *FSKeyinHandler) CheckSyntax() bool {
	if len(kh.options) != 0 {
		if kh.options == "ALL" {
			return nodeMgr.IsValidNodeName(kh.arguments)
		}
		return len(kh.options) <= 6 && len(kh.arguments) == 0
	}

	split := strings.Split(kh.arguments, ",")
	if len(split) < 1 {
		return false
	}

	for _, name := range split {
		if !nodeMgr.IsValidNodeName(strings.ToUpper(name)) {
			return false
		}
	}
	return true
}

func (kh *FSKeyinHandler) Invoke() {
	if !kh.threadStarted {
		go kh.thread()
	}
}

func (kh *FSKeyinHandler) IsDone() bool {
	return kh.threadStopped
}

func (kh *FSKeyinHandler) IsAllowed() bool {
	return true
}

func (kh *FSKeyinHandler) Dump(dest io.Writer, indent string) {
	_, _ = fmt.Fprintf(dest, "%vFS KEYIN ----------------------------------------------------\n", indent)

	_, _ = fmt.Fprintf(dest, "%v  threadStarted:  %v\n", indent, kh.threadStarted)
	_, _ = fmt.Fprintf(dest, "%v  threadStopped:  %v\n", indent, kh.threadStopped)
	_, _ = fmt.Fprintf(dest, "%v  terminateThread: %v\n", indent, kh.terminateThread)
}

func (kh *FSKeyinHandler) handleComponentList() {
	nm := kh.exec.GetNodeManager().(*nodeMgr.NodeManager)
	names := strings.Split(kh.arguments, ",")
	statStrings := make([]string, len(names))
	var err error
	for nx, name := range names {
		statStrings[nx], err = nm.GetNodeStatusStringForNode(strings.ToUpper(name))
		if err != nil {
			msg := fmt.Sprintf("FS KEYIN - %v DOES NOT EXIST, INPUT IGNORED", name)
			kh.exec.SendExecReadOnlyMessage(msg)
			return
		}
	}

	for sx := 0; sx < len(statStrings); {
		str := statStrings[sx]
		sx++
		if !strings.ContainsRune(str, '*') {
			if sx < len(statStrings) && !strings.ContainsRune(statStrings[sx], '*') {
				str = fmt.Sprintf("%-20s%s", str, statStrings[sx])
				sx++
			}
		}
		kh.exec.SendExecReadOnlyMessage(str)
	}
}

func (kh *FSKeyinHandler) handleOption() {
	switch kh.options {
	case "ALL":
		// TODO
	case "CM":
		// TODO
	case "DISK":
		// TODO
	case "FDISK":
		// TODO
	case "MS":
		// TODO
	case "PACK":
		// TODO
	case "RDISK":
		// TODO
	case "TAPE":
		// TODO
	}

	msg := fmt.Sprintf("FS KEYIN - %v OPTION DOES NOT EXIST, INPUT IGNORED", kh.options)
	kh.exec.SendExecReadOnlyMessage(msg)
}

func (kh *FSKeyinHandler) thread() {
	kh.threadStarted = true

	if len(kh.options) > 0 {
		kh.handleOption()
	} else {
		kh.handleComponentList()
	}

	kh.threadStopped = true
}
