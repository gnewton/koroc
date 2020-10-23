package main

import (
	"errors"
	"fmt"
	"strconv"
)

type SqlType int

const (
	UNKNOWN SqlType = iota
	Text
	Uint8  //TODO
	Int8   //TODO
	Uint16 //TODO
	Int16  //TODO
	Uint32
	Int32 //TODO
	Uint64
	Int64   //TODO
	Float32 //TODO
	Float64 //TODO
	Boolean
	Time //TODO
)

func (t SqlType) String() string {
	switch t {
	case Text:
		return "text"
	case Uint32:
		return "uint32"
	case Uint64:
		return "uint64"
	case Boolean:
		return "bool"
	}
	return "UNKNOWN TYPE"
}

func (t SqlType) ValueToString(field *Field, v interface{}) (string, error) {
	switch t {

	case Text:
		tv, ok := v.(string)
		if ok {
			return tv, nil
		} else {
			return "", errors.New("Value does not match field type: is not a string; field name:" + field.name + ":" + strconv.Itoa(field.positionInTable))
		}
	case Uint32:
		tv, ok := v.(uint32)
		if ok {
			return strconv.FormatUint(uint64(tv), 10), nil
		} else {
			return "", errors.New("Value does not match field type: is not a uint32" + Uint32.String())
		}
		//return "uint32"
	case Uint64:
		tv, ok := v.(uint64)
		if ok {
			return strconv.FormatUint(tv, 10), nil
		} else {
			return "", errors.New("Value does not match field type: is not a uint64" + Uint64.String())
		}
		//return "uint64"
	case Boolean:
		tv, ok := v.(bool)
		if ok {
			return strconv.FormatBool(tv), nil
		} else {
			return "", errors.New("Value does not match field type: is not a bool" + Boolean.String())
		}
		//return "bool"
	}
	return "", errors.New("Unknown type: " + t.String())
}

type Field struct {
	name            string
	typ             SqlType
	width           int
	positionInTable int
}

func (f *Field) CheckValueType(v interface{}) error {
	if v == nil {
		return errors.New("value is nil")
	}
	ok := true

	switch f.typ {
	case Text:
		_, ok = v.(string)
	case Uint32:
		_, ok = v.(uint32)
	case Uint64:
		_, ok = v.(uint64)
	case Boolean:
		_, ok = v.(bool)
	}

	if !ok {
		typ := ""
		switch tt := v.(type) {
		default:
			typ = fmt.Sprintf("%T", tt)
		}
		mes := fmt.Sprintf("Value does not match type value=%v matches=%v wantedType=%v actualType=%s", v, ok, f.typ, typ)
		return errors.New(mes)
	} else {
		return nil
	}

}
