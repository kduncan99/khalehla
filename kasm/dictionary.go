// Khalehla Project
// simple assembler
// Copyright © 2023 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"strings"
)

type DictionaryEntryType int

type DictionaryEntry interface {
	GetDictionaryEntryType() DictionaryEntryType
}

const (
	NodeDictionaryEntryType DictionaryEntryType = iota
	FunctionDictionaryEntryType
	ProcedureDictionaryEntryType
	ValueDictionaryEntryType
)

type Dictionary struct {
	parent  *Dictionary
	entries map[string]DictionaryEntry
}

func NewTopLevelDictionary() *Dictionary {
	return &Dictionary{
		parent:  nil,
		entries: make(map[string]DictionaryEntry),
	}
}

func NewSubLevelDictionary(parent *Dictionary) *Dictionary {
	return &Dictionary{
		parent:  parent,
		entries: make(map[string]DictionaryEntry),
	}
}

func (d *Dictionary) Establish(tag string, entry DictionaryEntry, level int) error {
	if level < 1 || d.parent == nil {
		upper := strings.ToUpper(tag)
		_, ok := d.entries[upper]
		if ok {
			return fmt.Errorf("duplicate label")
		}
		d.entries[upper] = entry
		return nil
	} else {
		return d.parent.Establish(tag, entry, level-1)
	}
}

func (d *Dictionary) EstablishFunction(tag string, function *Function, level int) error {
	fe := &FunctionDictionaryEntry{
		function: function,
	}
	return d.Establish(tag, fe, level)
}

func (d *Dictionary) EstablishProcedure(tag string, procedure *Procedure, level int) error {
	pe := &ProcedureDictionaryEntry{
		procedure: procedure,
	}
	return d.Establish(tag, pe, level)
}

func (d *Dictionary) EstablishValue(tag string, value Value, level int) error {
	ve := &ValueDictionaryEntry{
		value: value,
	}
	return d.Establish(tag, ve, level)
}

func (d *Dictionary) Lookup(tag string) (DictionaryEntry, error) {
	entry, ok := d.entries[strings.ToUpper(tag)]
	if ok {
		return entry, nil
	} else if d.parent != nil {
		return d.Lookup(tag)
	} else {
		return nil, fmt.Errorf("symbol not found")
	}
}

func (d *Dictionary) Remove(tag string) error {
	upper := strings.ToUpper(tag)
	_, ok := d.entries[upper]
	if ok {
		delete(d.entries, upper)
		return nil
	} else {
		return fmt.Errorf("symbol not found")
	}
}

// Node entry ----------------------------------------------------------------------------------------------------------

type NodeDictionaryEntry struct {
	entries map[uint64]DictionaryEntry
}

func (de *NodeDictionaryEntry) GetDictionaryEntryType() DictionaryEntryType {
	return NodeDictionaryEntryType
}

func (de *NodeDictionaryEntry) GetLeaf(selectors []uint64) (DictionaryEntry, error) {
	if len(selectors) == 0 {
		return de, nil
	} else {
		result, ok := de.entries[selectors[0]]
		if ok {
			if len(selectors) > 1 {
				if result.GetDictionaryEntryType() != NodeDictionaryEntryType {
					return nil, fmt.Errorf("too many selectors")
				} else {
					return result.(*NodeDictionaryEntry).GetLeaf(selectors[1:])
				}
			} else {
				return result, nil
			}
		} else {
			return nil, fmt.Errorf("selector not found")
		}
	}
}

// Function entry ------------------------------------------------------------------------------------------------------

type FunctionDictionaryEntry struct {
	function *Function
}

func (de *FunctionDictionaryEntry) GetDictionaryEntryType() DictionaryEntryType {
	return FunctionDictionaryEntryType
}

func (de *FunctionDictionaryEntry) GetFunction() *Function {
	return de.function
}

// Procedure entry -----------------------------------------------------------------------------------------------------

type ProcedureDictionaryEntry struct {
	procedure *Procedure
}

func (de *ProcedureDictionaryEntry) GetDictionaryEntryType() DictionaryEntryType {
	return ProcedureDictionaryEntryType
}

func (de *ProcedureDictionaryEntry) GetProcedure() *Procedure {
	return de.procedure
}

// Value entry ---------------------------------------------------------------------------------------------------------

type ValueDictionaryEntry struct {
	value Value
}

func (de *ValueDictionaryEntry) GetDictionaryEntryType() DictionaryEntryType {
	return ValueDictionaryEntryType
}

func (de *ValueDictionaryEntry) GetValue() Value {
	return de.value
}
