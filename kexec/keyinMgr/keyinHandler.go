// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package keyinMgr

import "time"

type KeyinHandler interface {
	Abort()
	CheckSyntax() bool
	GetCommand() string
	GetOptions() string
	GetArguments() string
	GetTimeFinished() time.Time
	Invoke()
	IsAllowed() bool
	IsDone() bool
}
