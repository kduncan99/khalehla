// Khalehla Project
// Copyright © 2023-2024 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package csi

// TODO ensure filespec is valid - parse it into a FileSpecification struct
// format:
//   [ [qualifier] '*' ] filename [cycle] [ '/' [read_key] [ '/' [write_key] ] ] ['.']
// cycle:
//   '(' [ '-' n1 ] | '0' | [ '+1' ] | n2 ')'
// n1: integer from 1 to 31
// n2: integer from 1 to 999

type FileSpecification struct {
	qualifier     *string // nil if no '*', empty if '*' with no qualifier
	filename      string
	absoluteCycle *int
	relativeCycle *int
	readKey       string
	writeKey      string
}

type parseIndex struct {
	text  string
	index int
}

func (pi *parseIndex) parseIndexAtEnd() bool {
	return pi.index >= len(pi.text)
}

func parseQualifierAndFilename(pi *parseIndex) (qualifier *string, filename string, ok bool) {
	// TODO
	return nil, "", false
}

func parseCycle(pi *parseIndex) (absoluteCycle *int, relativeCycle *int, ok bool) {
	// TODO
	return nil, nil, true
}

func parseKeys(pi *parseIndex) (readkey string, writekey string, ok bool) {
	// TODO
	return "", "", true
}

func ParseFileSpecification(input string) (*FileSpecification, bool) {
	pi := &parseIndex{
		text:  input,
		index: 0,
	}

	qualifier, filename, ok := parseQualifierAndFilename(pi)
	if !ok {
		return nil, ok
	}

	abs, rel, ok := parseCycle(pi)
	if !ok {
		return nil, ok
	}

	rkey, wkey, ok := parseKeys(pi)
	if !ok {
		return nil, ok
	}

	// TODO check terminating '.', but nothing else is allowed

	return &FileSpecification{
		qualifier:     qualifier,
		filename:      filename,
		absoluteCycle: abs,
		relativeCycle: rel,
		readKey:       rkey,
		writeKey:      wkey,
	}, true
}
