// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package wire

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

// LocalInstrument is the LocalInstrument of the wire
type LocalInstrument struct {
	// tag
	tag string
	// LocalInstrumentCode is local instrument code
	LocalInstrumentCode string `json:"LocalInstrument,omitempty"`
	// ProprietaryCode is proprietary code
	ProprietaryCode string `json:"proprietaryCode,omitempty"`

	// validator is composed for data validation
	validator
	// converters is composed for WIRE to GoLang Converters
	converters
}

// NewLocalInstrument returns a new LocalInstrument
func NewLocalInstrument() *LocalInstrument {
	li := &LocalInstrument{
		tag: TagLocalInstrument,
	}
	return li
}

// Parse takes the input string and parses the LocalInstrument values
//
// Parse provides no guarantee about all fields being filled in. Callers should make a Validate() call to confirm
// successful parsing and data validity.
func (li *LocalInstrument) Parse(record string) error {
	if utf8.RuneCountInString(record) < 6 {
		return NewTagMinLengthErr(6, len(record))
	}

	li.tag = record[:6]
	length := 6

	value, read, err := li.parseVariableStringField(record[length:], 4)
	if err != nil {
		return fieldError("SwiftFieldTag", err)
	}
	li.LocalInstrumentCode = value
	length += read

	value, read, err = li.parseVariableStringField(record[length:], 35)
	if err != nil {
		return fieldError("ProprietaryCode", err)
	}
	li.ProprietaryCode = value
	length += read

	if len(record) != length {
		return NewTagMaxLengthErr()
	}

	return nil
}

func (li *LocalInstrument) UnmarshalJSON(data []byte) error {
	type Alias LocalInstrument
	aux := struct {
		*Alias
	}{
		(*Alias)(li),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	li.tag = TagLocalInstrument
	return nil
}

// String writes LocalInstrument
func (li *LocalInstrument) String(options ...bool) string {
	var buf strings.Builder
	buf.Grow(45)

	buf.WriteString(li.tag)
	buf.WriteString(li.LocalInstrumentCodeField(options...))
	buf.WriteString(li.ProprietaryCodeField(options...))

	if li.parseFirstOption(options) {
		return li.stripDelimiters(buf.String())
	} else {
		return buf.String()
	}
}

// Validate performs WIRE format rule checks on LocalInstrument and returns an error if not Validated
// The first error encountered is returned and stops that parsing.
func (li *LocalInstrument) Validate() error {
	if err := li.fieldInclusion(); err != nil {
		return err
	}
	if li.tag != TagLocalInstrument {
		return fieldError("tag", ErrValidTagForType, li.tag)
	}
	if err := li.isLocalInstrumentCode(li.LocalInstrumentCode); err != nil {
		return fieldError("LocalInstrumentCode", err, li.LocalInstrumentCode)
	}
	if err := li.isAlphanumeric(li.ProprietaryCode); err != nil {
		return fieldError("ProprietaryCode", err, li.ProprietaryCode)
	}
	return nil
}

// fieldInclusion validate mandatory fields. If fields are
// invalid the WIRE will return an error.
// ProprietaryCode is only allowed if LocalInstrument Code is PROP
func (li *LocalInstrument) fieldInclusion() error {
	if li.LocalInstrumentCode != ProprietaryLocalInstrumentCode && li.ProprietaryCode != "" {
		return fieldError("ProprietaryCode", ErrInvalidProperty, li.ProprietaryCode)
	}
	return nil
}

// LocalInstrumentCodeField gets a string of LocalInstrumentCode field
func (li *LocalInstrument) LocalInstrumentCodeField(options ...bool) string {
	return li.alphaVariableField(li.LocalInstrumentCode, 4, li.parseFirstOption(options))
}

// ProprietaryCodeField gets a string of ProprietaryCode field
func (li *LocalInstrument) ProprietaryCodeField(options ...bool) string {
	return li.alphaVariableField(li.ProprietaryCode, 35, li.parseFirstOption(options))
}
