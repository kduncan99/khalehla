// khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import "time"

type IKeyinHandler interface {
	Abort()
	CheckSyntax() bool
	GetArguments() string
	GetCommand() string
	GetOptions() string
	GetHelp() []string
	GetTimeFinished() time.Time
	Invoke()
	IsAllowed() bool
	IsDone() bool
}
