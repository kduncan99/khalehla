// khalehla Project
// simple assembler
// Copyright © 2023-2025 by Kurt Duncan, BearSnake LLC
// All Rights Reserved

package kasm

import (
	"fmt"

	"khalehla/common"
	pkg2 "khalehla/old/pkg"
)

type NodeValue struct {
	entries map[uint64]Value
}

var invalidValueType = fmt.Errorf("invalid value type for selector")
var compositeError = fmt.Errorf("composite values cannot be used as selectors")
var invalidValue = fmt.Errorf("integer value invalid for selector")
var undefinedOffsetError = fmt.Errorf("cannot use values with undefined offsets as selectors")

func NewNodeValue() *NodeValue {
	return &NodeValue{
		entries: make(map[uint64]Value),
	}
}

func NewNodeWithLeaf(selectors []Value, leaf Value) (*NodeValue, error) {
	if len(selectors) == 0 {
		return nil, fmt.Errorf("no selectors provided")
	}

	nv := NewNodeValue()
	current := nv
	sel := selectors
	for len(sel) > 0 {
		subSel := sel[1:]
		ix, err := getIndexFromSelectorValue(sel[0])
		if err != nil {
			return nil, err
		}

		if len(subSel) > 0 {
			current.entries[ix] = NewNodeValue()
			sel = subSel
		} else {
			current.entries[ix] = leaf
		}
	}

	return nv, nil
}

func (nv *NodeValue) GetValueType() ValueType {
	return NodeValueType
}

func (nv *NodeValue) Evaluate(ec *ExpressionContext) error {
	selectors, err := ec.PopVariableParameterList()
	if err != nil {
		return err
	}

	value, err := nv.eval(selectors)
	if err != nil {
		return err
	}

	ec.PushValue(value)
	return nil
}

func (nv *NodeValue) Merge(selectors []Value, value Value) error {
	if len(selectors) == 0 {
		return fmt.Errorf("attempt to set a node to a leaf value")
	}

	selector := selectors[0]
	subSelectors := selectors[1:]
	index, err := getIndexFromSelectorValue(selector)
	if err != nil {
		return err
	}

	val, err := nv.getValueAt(selector)
	if err != nil {
		return err
	}

	if val == nil {
		// nothing at the given index
		if len(subSelectors) == 0 {
			nv.entries[index] = value
		} else {
			nv.entries[index], err = NewNodeWithLeaf(subSelectors, value)
			if err != nil {
				return err
			}
		}
	} else if val.GetValueType() != NodeValueType {
		//	Non-node at the given index
		if len(selectors) > 0 {
			return fmt.Errorf("attempt to override a leaf with a node")
		} else {
			nv.entries[index] = value
		}
	} else {
		//	Node at the given index
		if len(selectors) == 0 {
			return fmt.Errorf("attempt to override a node with a leaf")
		} else {
			return val.(*NodeValue).Merge(selectors[1:], value)
		}
	}

	return nil
}

func (nv *NodeValue) eval(selectors []Value) (Value, error) {
	if len(selectors) == 0 {
		return NewSimpleIntegerValue(uint64(len(nv.entries))), nil
	}

	value, err := nv.getValueAt(selectors[0])
	if err != nil {
		return nil, err
	} else if value == nil {
		return nil, fmt.Errorf("cannot find selector %v in node", selectors[0])
	}

	subSelectors := selectors[1:]
	if value.GetValueType() == NodeValueType {
		nv := value.(*NodeValue)
		return nv.eval(subSelectors)
	}

	if len(subSelectors) > 0 {
		return nil, fmt.Errorf("attempt to use selectors on a non-node value")
	}

	return value, nil
}

func getIndexFromSelectorValue(selector Value) (uint64, error) {
	if selector.GetValueType() != IntegerValueType {
		return 0, invalidValueType
	}

	iVal := selector.(*IntegerValue)
	if len(iVal.components) != 1 {
		return 0, compositeError
	}

	comp := iVal.components[0]
	if len(comp.offsets) > 0 {
		return 0, undefinedOffsetError
	}

	if (comp.value&pkg2.NegativeZero != comp.value) || common.IsNegative(comp.value) {
		return 0, invalidValue
	}

	return comp.value, nil
}

func (nv *NodeValue) getValueAt(selector Value) (Value, error) {
	index, err := getIndexFromSelectorValue(selector)
	if err != nil {
		return nil, err
	}

	val, ok := nv.entries[index]
	if !ok {
		return nil, nil
	} else {
		return val, nil
	}
}
