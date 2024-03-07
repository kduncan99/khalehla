// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package facilitiesMgr

import (
	"fmt"
	"khalehla/kexec"
	"khalehla/kexec/mfdMgr"
)

type FileSpecification struct {
	Qualifier     string
	HasAsterisk   bool // if true but qualifier is empty, use the implied qualifier
	Filename      string
	FileCycleSpec *mfdMgr.FileCycleSpecification
	ReadKey       string
	WriteKey      string
}

func (fs *FileSpecification) parseQualifierAndFilename(p *kexec.Parser) (fsCode FacStatusCode, ok bool) {
	fsCode = 0
	ok = true

	// There must be at least a filename, else we are in error
	fs.Qualifier = ""
	fs.HasAsterisk = false
	fs.Filename = ""

	for !p.IsAtEnd() {
		ch, _ := p.ParseNextCharacter()
		if !fs.HasAsterisk && kexec.IsValidQualifierChar(ch) {
			if len(fs.Qualifier) == 12 {
				return FacStatusSyntaxErrorInImage, false
			} else {
				fs.Qualifier += string(rune(ch))
				continue
			}
		}

		if fs.HasAsterisk && kexec.IsValidFilenameChar(ch) {
			if len(fs.Filename) == 12 {
				return FacStatusSyntaxErrorInImage, false
			} else {
				fs.Filename += string(rune(ch))
				continue
			}
		}

		if ch == '*' {
			if fs.HasAsterisk {
				return FacStatusSyntaxErrorInImage, false
			} else {
				fs.HasAsterisk = true
			}
		}
	}

	if len(fs.Filename) == 0 {
		return FacStatusFilenameIsRequired, false
	}

	return
}

func (fs *FileSpecification) parseAbsoluteCycle(p *kexec.Parser) (found bool, fsCode FacStatusCode, ok bool) {
	found = false
	fsCode = 0
	ok = true

	if !p.ParseSpecificCharacter('(') {
		return
	}

	found = true

	p.SkipSpaces()
	value, found, err := p.ParseUnsignedInteger()
	if err != nil || !found {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	if value < 1 || value > 999 {
		fsCode = FacStatusIllegalValueForFCycle
		ok = false
		return
	}

	p.SkipSpaces()
	if !p.ParseSpecificCharacter(')') {
		fsCode = FacStatusIllegalValueForFCycle
		ok = false
		return
	}

	fs.FileCycleSpec = mfdMgr.NewAbsoluteFileCycleSpecification(uint(value))
	return
}

func (fs *FileSpecification) parseRelativeCycle(p *kexec.Parser) (found bool, fsCode FacStatusCode, ok bool) {
	found = false
	fsCode = FacStatusComplete

	p.MarkPosition()
	if !p.ParseSpecificCharacter('(') {
		return
	}

	p.SkipSpaces()
	var pos bool
	var neg bool
	pos = p.ParseSpecificCharacter('+')
	if !pos {
		neg = p.ParseSpecificCharacter('-')
	}

	if !pos && !neg {
		p.ResetPosition()
		return
	}

	found = true

	p.SkipSpaces()
	value, found, err := p.ParseUnsignedInteger()
	if err != nil || !found {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	p.SkipSpaces()
	if !p.ParseSpecificCharacter(')') {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	if pos {
		if value == 1 {
			fs.FileCycleSpec = mfdMgr.NewRelativeFileCycleSpecification(int(value))
			found = true
			return
		} else {
			fsCode = FacStatusIllegalValueForFCycle
			ok = false
			return
		}
	} else {
		if value < 1 || value > 31 {
			fsCode = FacStatusIllegalValueForFCycle
			ok = false
			return
		} else {
			fs.FileCycleSpec = mfdMgr.NewRelativeFileCycleSpecification(int(value) * -1)
			found = true
			return
		}
	}
}

func (fs *FileSpecification) parseCycle(p *kexec.Parser) (fsCode FacStatusCode, ok bool) {
	found, fsCode, ok := fs.parseRelativeCycle(p)
	if !ok {
		return
	} else if !found {
		_, fsCode, ok = fs.parseAbsoluteCycle(p)
	}

	return
}

func (fs *FileSpecification) parseKey(p *kexec.Parser) (key string, err error) {
	key, _ = p.ParseUntil(",./; ")
	if !kexec.IsValidReadWriteKey(key) {
		err = fmt.Errorf("invalid key")
	}
	return
}

func (fs *FileSpecification) parseKeys(p *kexec.Parser) (fsCode FacStatusCode, ok bool) {
	fsCode = 0
	ok = true

	var err error
	fs.ReadKey, err = fs.parseKey(p)
	if err != nil {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	fs.WriteKey, err = fs.parseKey(p)
	if err != nil {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	return
}

// ParseFileSpecification parses the given input string in an attempt to decode the
// qualifier, file, cycle, read key, and write key subfields.
// format:
//
//	[ [qualifier] '*' ] filename [cycle] [ '/' [read_key] [ '/' [write_key] ] ] ['.']
//
// cycle:
//
//	'(' [ '-' n1 ] | '0' | [ '+1' ] | n2 ')'
//
// n1: integer from 1 to 31
// n2: integer from 1 to 999
// If the input is empty, we return nil in FileSpecification and ok == true.
// If successful, we return a pointer to the FileSpecification in fs, with ok == true.
// If we find something, but encounter an error during parsing, we return ok == false and something descriptive in code.
func ParseFileSpecification(p *kexec.Parser) (fs *FileSpecification, fsCode FacStatusCode, ok bool) {
	fs = nil
	fsCode = 0
	ok = true

	if p.IsAtEnd() {
		return
	}

	fs = &FileSpecification{}
	fsCode, ok = fs.parseQualifierAndFilename(p)
	if !ok {
		return
	}

	fsCode, ok = fs.parseCycle(p)
	if !ok {
		return
	}

	fsCode, ok = fs.parseKeys(p)
	if !ok {
		return
	}

	// eat terminating '.' if it exists
	p.ParseSpecificCharacter('.')

	return
}
