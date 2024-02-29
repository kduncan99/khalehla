// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kexec

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

func (fs *FileSpecification) parseQualifierAndFilename(p *Parser) (fsCode FacStatusCode, ok bool) {
	fsCode = 0
	ok = true

	// There must be at least a filename, else we are in error
	fs.Qualifier = ""
	fs.HasAsterisk = false
	fs.Filename = ""

	for !p.IsAtEnd() {
		ch, _ := p.ParseNextCharacter()
		if !fs.HasAsterisk && IsValidQualifierChar(ch) {
			if len(fs.Qualifier) == 12 {
				return FacStatusSyntaxErrorInImage, false
			} else {
				fs.Qualifier += string(rune(ch))
				continue
			}
		}

		if fs.HasAsterisk && IsValidFilenameChar(ch) {
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

func (fs *FileSpecification) parseAbsoluteCycle(p *Parser) (found bool, fsCode FacStatusCode, ok bool) {
	found = false
	fsCode = 0
	ok = true

	if !p.ParseSpecificCharacter('(') {
		return
	}

	found = true

	p.SkipSpaces()
	value, err := p.ParseUnsignedInteger()
	if err != nil {
		fsCode = FacStatusSyntaxErrorInImage
		ok = false
		return
	}

	if *value < 1 || *value > 999 {
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

	abs := uint(*value)
	fs.AbsoluteCycle = &abs
	return
}

func (fs *FileSpecification) parseRelativeCycle(p *Parser) (found bool, fsCode FacStatusCode, ok bool) {
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
	value, err := p.ParseUnsignedInteger()
	if err != nil {
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
		if *value == 1 {
			iVal := int(*value)
			fs.RelativeCycle = &iVal
			found = true
			return
		} else {
			fsCode = FacStatusIllegalValueForFCycle
			ok = false
			return
		}
	} else {
		if *value < 1 || *value > 31 {
			fsCode = FacStatusIllegalValueForFCycle
			ok = false
			return
		} else {
			iVal := int(*value) * -1
			fs.RelativeCycle = &iVal
			found = true
			return
		}
	}
}

func (fs *FileSpecification) parseCycle(p *Parser) (fsCode FacStatusCode, ok bool) {
	found, fsCode, ok := fs.parseRelativeCycle(p)
	if !ok {
		return
	} else if !found {
		_, fsCode, ok = fs.parseAbsoluteCycle(p)
	}

	return
}

func (fs *FileSpecification) parseKeys(p *Parser) (fsCode FacStatusCode, ok bool) {
	fsCode = 0
	ok = true

	// TODO
	return
}

// ParseFileSpecification parses the given input string in an attempt to decode the
// qualifier, file, cycle, read key, and write key subfields.
// If the input is empty, we return nil in FileSpecification and ok == true.
// If successful, we return a pointer to the FileSpecification in fs, with ok == true.
// If we find something, but encounter an error during parsing, we return ok == false and something descriptive in code.
func ParseFileSpecification(p *Parser) (fs *FileSpecification, fsCode FacStatusCode, ok bool) {
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
