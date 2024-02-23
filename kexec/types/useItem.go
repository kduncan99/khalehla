// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package types

type UseItem struct {
	internalFilename    string
	referencedQualifier string
	impliedQualifier    bool // asterisk was used without a qualifier - use implied qualifier
	referencedFilename  string
	releaseFlag         bool
}
