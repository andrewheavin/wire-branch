// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package wire

import (
	"encoding/json"
	"strings"
	"unicode/utf8"
)

// Originator is the originator of the wire
type Originator struct {
	// tag
	tag string
	// Personal
	Personal Personal `json:"personal,omitempty"`

	// validator is composed for data validation
	validator
	// converters is composed for WIRE to GoLang Converters
	converters
}

// NewOriginator returns a new Originator
func NewOriginator() *Originator {
	o := &Originator{
		tag: TagOriginator,
	}
	return o
}

// Parse takes the input string and parses the Originator values
//
// Parse provides no guarantee about all fields being filled in. Callers should make a Validate() call to confirm
// successful parsing and data validity.
func (o *Originator) Parse(record string) error {
	if utf8.RuneCountInString(record) < 9 {
		return NewTagMinLengthErr(9, len(record))
	}

	o.tag = record[:6]
	o.Personal.IdentificationCode = o.parseStringField(record[6:7])
	length := 7

	value, read, err := o.parseVariableStringField(record[length:], 34)
	if err != nil {
		return fieldError("Identifier", err)
	}
	o.Personal.Identifier = value
	length += read

	value, read, err = o.parseVariableStringField(record[length:], 35)
	if err != nil {
		return fieldError("Name", err)
	}
	o.Personal.Name = value
	length += read

	value, read, err = o.parseVariableStringField(record[length:], 35)
	if err != nil {
		return fieldError("AddressLineOne", err)
	}
	o.Personal.Address.AddressLineOne = value
	length += read

	value, read, err = o.parseVariableStringField(record[length:], 35)
	if err != nil {
		return fieldError("AddressLineTwo", err)
	}
	o.Personal.Address.AddressLineTwo = value
	length += read

	value, read, err = o.parseVariableStringField(record[length:], 35)
	if err != nil {
		return fieldError("AddressLineThree", err)
	}
	o.Personal.Address.AddressLineThree = value
	length += read

	if len(record) != length {
		return NewTagMaxLengthErr()
	}

	return nil
}

func (o *Originator) UnmarshalJSON(data []byte) error {
	type Alias Originator
	aux := struct {
		*Alias
	}{
		(*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	o.tag = TagOriginator
	return nil
}

// String writes Originator
func (o *Originator) String(options ...bool) string {
	var buf strings.Builder
	buf.Grow(181)

	buf.WriteString(o.tag)
	buf.WriteString(o.IdentificationCodeField())
	buf.WriteString(o.IdentifierField(options...))
	buf.WriteString(o.NameField(options...))
	buf.WriteString(o.AddressLineOneField(options...))
	buf.WriteString(o.AddressLineTwoField(options...))
	buf.WriteString(o.AddressLineThreeField(options...))

	if o.parseFirstOption(options) {
		return o.stripDelimiters(buf.String())
	} else {
		return buf.String()
	}
}

// Validate performs WIRE format rule checks on Originator and returns an error if not Validated
// The first error encountered is returned and stops that parsing.
func (o *Originator) Validate() error {
	if err := o.fieldInclusion(); err != nil {
		return err
	}
	if o.tag != TagOriginator {
		return fieldError("tag", ErrValidTagForType, o.tag)
	}
	// Can be any Identification Code
	if err := o.isIdentificationCode(o.Personal.IdentificationCode); err != nil {
		return fieldError("IdentificationCode", err, o.Personal.IdentificationCode)
	}
	if err := o.isAlphanumeric(o.Personal.Identifier); err != nil {
		return fieldError("Identifier", err, o.Personal.Identifier)
	}
	if err := o.isAlphanumeric(o.Personal.Name); err != nil {
		return fieldError("Name", err, o.Personal.Name)
	}
	if err := o.isAlphanumeric(o.Personal.Address.AddressLineOne); err != nil {
		return fieldError("AddressLineOne", err, o.Personal.Address.AddressLineOne)
	}
	if err := o.isAlphanumeric(o.Personal.Address.AddressLineTwo); err != nil {
		return fieldError("AddressLineTwo", err, o.Personal.Address.AddressLineTwo)
	}
	if err := o.isAlphanumeric(o.Personal.Address.AddressLineThree); err != nil {
		return fieldError("AddressLineThree", err, o.Personal.Address.AddressLineThree)
	}
	return nil
}

// fieldInclusion validate mandatory fields. If fields are
// invalid the WIRE will return an error.
func (o *Originator) fieldInclusion() error {
	if o.Personal.IdentificationCode != "" && o.Personal.Identifier == "" {
		return fieldError("Identifier", ErrFieldRequired)
	}
	if o.Personal.IdentificationCode == "" && o.Personal.Identifier != "" {
		return fieldError("IdentificationCode", ErrFieldRequired)
	}
	return nil
}

// IdentificationCodeField gets a string of the IdentificationCode field
func (o *Originator) IdentificationCodeField() string {
	return o.alphaField(o.Personal.IdentificationCode, 1)
}

// IdentifierField gets a string of the Identifier field
func (o *Originator) IdentifierField(options ...bool) string {
	return o.alphaVariableField(o.Personal.Identifier, 34, o.parseFirstOption(options))
}

// NameField gets a string of the Name field
func (o *Originator) NameField(options ...bool) string {
	return o.alphaVariableField(o.Personal.Name, 35, o.parseFirstOption(options))
}

// AddressLineOneField gets a string of AddressLineOne field
func (o *Originator) AddressLineOneField(options ...bool) string {
	return o.alphaVariableField(o.Personal.Address.AddressLineOne, 35, o.parseFirstOption(options))
}

// AddressLineTwoField gets a string of AddressLineTwo field
func (o *Originator) AddressLineTwoField(options ...bool) string {
	return o.alphaVariableField(o.Personal.Address.AddressLineTwo, 35, o.parseFirstOption(options))
}

// AddressLineThreeField gets a string of AddressLineThree field
func (o *Originator) AddressLineThreeField(options ...bool) string {
	return o.alphaVariableField(o.Personal.Address.AddressLineThree, 35, o.parseFirstOption(options))
}
