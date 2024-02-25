// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

import (
	"fmt"
)

// ensure filespec is valid - parse it into a FileSpecification struct
// format:
//   [ [qualifier] '*' ] filename [cycle] [ '/' [read_key] [ '/' [write_key] ] ] ['.']
// cycle:
//   '(' [ '-' n1 ] | '0' | [ '+1' ] | n2 ')'
// n1: integer from 1 to 31
// n2: integer from 1 to 999

type FileSpecification struct {
	Qualifier     string
	HasAsterisk   bool // if true but qualifier is empty, use the implied qualifier
	Filename      string
	AbsoluteCycle *uint
	RelativeCycle *int
	ReadKey       string
	WriteKey      string
}

func (fs *FileSpecification) parseQualifierAndFilename(p *Parser) error {
	// There must be at least a filename, else we are in error
	fs.Qualifier = ""
	fs.HasAsterisk = false
	fs.Filename = ""

	for !p.IsAtEnd() {
		ch, _ := p.ParseNextCharacter()
		if !fs.HasAsterisk && IsValidQualifierChar(ch) {
			if len(fs.Qualifier) == 12 {
				return fmt.Errorf("invalid qualifier") // error
			} else {
				fs.Qualifier += string(rune(ch))
				continue
			}
		}

		if fs.HasAsterisk && IsValidFilenameChar(ch) {
			if len(fs.Filename) == 12 {
				return fmt.Errorf("invalid filename")
			} else {
				fs.Filename += string(rune(ch))
				continue
			}
		}

		if ch == '*' {
			if fs.HasAsterisk {
				return fmt.Errorf("invalid format")
			} else {
				fs.HasAsterisk = true
			}
		}
	}

	if len(fs.Filename) == 0 {
		return fmt.Errorf("no filename")
	}

	return nil
}

func (fs *FileSpecification) parseAbsoluteCycle(p *Parser) (found bool, err error) {
	found = false
	err = nil

	if !p.ParseSpecificCharacter('(') {
		return
	}

	p.SkipSpaces()
	value, err := p.ParseUnsignedInteger()
	if err != nil {
		return
	}

	if *value < 1 || *value > 999 {
		err = fmt.Errorf("invalid absolute file cycle")
		return
	}

	p.SkipSpaces()
	if !p.ParseSpecificCharacter(')') {
		err = fmt.Errorf("syntax error in file cycle")
		return
	}

	abs := uint(*value)
	fs.AbsoluteCycle = &abs
	found = true
	return
}

func (fs *FileSpecification) parseRelativeCycle(p *Parser) (found bool, err error) {
	found = false
	err = nil

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

	p.SkipSpaces()
	value, err := p.ParseUnsignedInteger()
	if err != nil {
		return
	}

	p.SkipSpaces()
	if !p.ParseSpecificCharacter(')') {
		err = fmt.Errorf("syntax error in file cycle")
		return
	}

	if pos {
		if *value == 1 {
			iVal := int(*value)
			fs.RelativeCycle = &iVal
			found = true
			return
		} else {
			err = fmt.Errorf("illegal relative file cycle")
			return
		}
	} else {
		if *value < 1 || *value > 31 {
			err = fmt.Errorf("illegal relative file cycle")
			return
		} else {
			iVal := int(*value) * -1
			fs.RelativeCycle = &iVal
			found = true
			return
		}
	}
}

func (fs *FileSpecification) parseCycle(p *Parser) error {
	found, err := fs.parseRelativeCycle(p)
	if err != nil {
		return err
	}

	if !found {
		_, err := fs.parseAbsoluteCycle(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileSpecification) parseKeys(p *Parser) error {
	// TODO
	return nil
}

func ParseFileSpecification(input string) (fs *FileSpecification, found bool, err error) {
	// I think the input should have no embedded spaces...
	// TODO if the above is true, we might be able to simplify some code...?
	fs = &FileSpecification{}
	found = false
	err = nil

	p := NewParser(input)
	if p.IsAtEnd() {
		return
	}

	found = true

	err = fs.parseQualifierAndFilename(p)
	if err != nil {
		return
	}

	err = fs.parseCycle(p)
	if err != nil {
		return
	}

	err = fs.parseKeys(p)
	if err != nil {
		return
	}

	// TODO check terminating '.', but nothing else is allowed

	found = true
	return
}
