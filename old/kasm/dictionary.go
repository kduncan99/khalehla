// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"
	"strings"
)

type Dictionary struct {
	parent  *Dictionary
	entries map[string]Value
}

var noSelectors []Value

func NewTopLevelDictionary() *Dictionary {
	d := &Dictionary{
		parent:  nil,
		entries: make(map[string]Value),
	}

	for token, fn := range Functions {
		_ = d.Establish(token, noSelectors, 0, fn)
	}

	return d
}

func NewSubLevelDictionary(parent *Dictionary) *Dictionary {
	return &Dictionary{
		parent:  parent,
		entries: make(map[string]Value),
	}
}

func (d *Dictionary) Establish(tag string, selectors []Value, level int, entry Value) error {
	if level > 0 && d.parent != nil {
		return d.parent.Establish(tag, selectors, level-1, entry)
	}

	tagUpper := strings.ToUpper(tag)
	entry, ok := d.entries[tagUpper]
	if ok {
		//	the tag exists in the dictionary
		if entry.GetValueType() != NodeValueType || len(selectors) == 0 {
			return fmt.Errorf("duplicate symbol")
		}

		return entry.(*NodeValue).Merge(selectors, entry)
	} else {
		//	the tag does not exist in the dictionary
		if len(selectors) == 0 {
			//	no selectors - create an entry in the dictionary and we're done.
			d.entries[tagUpper] = entry
			return nil
		} else {
			//  create a new node in the dictionary
			node, err := NewNodeWithLeaf(selectors, entry)
			if node != nil {
				d.entries[tagUpper] = node
			}
			return err
		}
	}
}

func (d *Dictionary) Lookup(tag string) (Value, error) {
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
