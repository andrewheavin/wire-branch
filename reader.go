// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package wire

import (
	"bufio"
	"github.com/moov-io/base"
	"io"
)

// Reader reads records from a ACH-encoded file.
type Reader struct {
	// r handles the IO.Reader sent to be parser.
	scanner *bufio.Scanner
	// file is ach.file model being built as r is parsed.
	File File
	// line is the current line being parsed from the input r
	line string
	// ToDo:  Do we need a current FEDWireMessage, just use FEDWireMessage
	// currentFEDWireMessage is the current FEDWireMessage being parsed
	currentFEDWireMessage FEDWireMessage
	// lineNum is the line number of the file being parsed
	lineNum int
	// tagName holds the current tag name being parsed.
	tagName string
	// errors holds each error encountered when attempting to parse the file
	errors base.ErrorList
}

// error returns a new ParseError based on err
func (r *Reader) parseError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(*base.ParseError); ok {
		return err
	}
	return &base.ParseError{
		Line:   r.lineNum,
		Record: r.tagName,
		Err:    err,
	}
}

// NewReader returns a new ACH Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(r),
	}
}

// addCurrentFEDWireMessage creates the current FEDWireMessage for the file being read. A successful
// current FEDWireMessage will be added to r.File once parsed.
func (r *Reader) addCurrentFEDWireMessage(fwm FEDWireMessage) {
	r.currentFEDWireMessage = FEDWireMessage{}
}

// Read reads each line of the FED Wire file and defines which parser to use based
// on the first character of each line. It also enforces FED Wire formatting rules and returns
// the appropriate error if issues are found.
func (r *Reader) Read() (File, error) {
	r.lineNum = 0
	// read through the entire file
	for r.scanner.Scan() {
		line := r.scanner.Text()
		r.lineNum++
		// ToDo: File length Check?
		r.line = line
		if err := r.parseLine(); err != nil {
			r.errors.Add(err)
		}
	}

	r.File.AddFEDWireMessage(r.currentFEDWireMessage)
	r.currentFEDWireMessage = NewFEDWireMessage()

	if r.errors.Empty() {
		return r.File, nil
	}
	return r.File, r.errors
}

func (r *Reader) parseLine() error {
	switch r.line[:6] {
	case TagSenderSupplied:
		if err := r.parseSenderSupplied(); err != nil {
			return err
		}
	case TagTypeSubType:
		if err := r.parseTypeSubType(); err != nil {
			return err
		}
	case TagInputMessageAccountabilityData:
		if err := r.parseInputMessageAccountabilityData(); err != nil {
			return err
		}
	case TagAmount:
		if err := r.parseAmount(); err != nil {
			return err
		}
	case TagSenderDepositoryInstitution:
		if err := r.parseSenderDepositoryInstitution(); err != nil {
			return err
		}
	case TagReceiverDepositoryInstitution:
		if err := r.parseReceiverDepositoryInstitution(); err != nil {
			return err
		}
	case TagBusinessFunctionCode:
		if err := r.parseBusinessFunctionCode(); err != nil {
			return err
		}
	case TagSenderReference:
		if err := r.parseSenderReference(); err != nil {
			return err
		}
	case TagPreviousMessageIdentifier:
		if err := r.parsePreviousMessageIdentifier(); err != nil {
			return err
		}
	case TagLocalInstrument:
		if err := r.parseLocalInstrument(); err != nil {
			return err
		}
	case TagPaymentNotification:
		if err := r.parsePaymentNotification(); err != nil {
			return err
		}
	case TagCharges:
		if err := r.parseCharges(); err != nil {
			return err
		}
	case TagInstructedAmount:
		if err := r.parseInstructedAmount(); err != nil {
			return err
		}
	case TagExchangeRate:
		if err := r.parseExchangeRate(); err != nil {
			return err
		}

	case TagBeneficiaryIntermediaryFI:
		if err := r.parseBeneficiaryIntermediaryFI(); err != nil {
			return err
		}
	case TagBeneficiaryFI:
		if err := r.parseBeneficiaryFI(); err != nil {
			return err
		}
	case TagBeneficiary:
		if err := r.parseBeneficiary(); err != nil {
			return err
		}
	case TagBeneficiaryReference:
		if err := r.parseBeneficiaryReference(); err != nil {
			return err
		}
	case TagAccountDebitedDrawdown:
		if err := r.parseAccountDebitedDrawdown(); err != nil {
			return err
		}
	case TagOriginator:
		if err := r.parseOriginator(); err != nil {
			return err
		}
	case TagOriginatorOptionF:
		if err := r.parseOriginatorOptionF(); err != nil {
			return err
		}
	case TagOriginatorFI:
		if err := r.parseOriginatorFI(); err != nil {
			return err
		}
	case TagInstructingFI:
		if err := r.parseInstructingFI(); err != nil {
			return err
		}
	case TagAccountCreditedDrawdown:
		if err := r.parseAccountCreditedDrawdown(); err != nil {
			return err
		}
	case TagOriginatorToBeneficiary:
		if err := r.parseOriginatorToBeneficiary(); err != nil {
			return err
		}
	case TagFIReceiverFI:
		if err := r.parseFIReceiverFI(); err != nil {
			return err
		}
	case TagFIDrawdownDebitAccountAdvice:
		if err := r.parseFIDrawdownDebitAccountAdvice(); err != nil {
			return err
		}
	case TagFIIntermediaryFI:
		if err := r.parseFIIntermediaryFI(); err != nil {
			return err
		}
	case TagFIIntermediaryFIAdvice:
		if err := r.parseFIIntermediaryFIAdvice(); err != nil {
			return err
		}
	case TagFIBeneficiaryFI:
		if err := r.parseFIBeneficiaryFI(); err != nil {
			return err
		}
	case TagFIBeneficiaryFIAdvice:
		if err := r.parseFIBeneficiaryFIAdvice(); err != nil {
			return err
		}
	case TagFIBeneficiary:
		if err := r.parseFIBeneficiary(); err != nil {
			return err
		}
	case TagFIBeneficiaryAdvice:
		if err := r.parseFIBeneficiaryAdvice(); err != nil {
			return err
		}
	case TagFIPaymentMethodToBeneficiary:
		if err := r.parseFIPaymentMethodToBeneficiary(); err != nil {
			return err
		}
	case TagFIAdditionalFIToFI:
		if err := r.parseFIAdditionalFiToFi(); err != nil {
			return err
		}
	case TagCurrencyInstructedAmount:
		if err := r.parseCurrencyInstructedAmount(); err != nil {
			return err
		}
	case TagOrderingCustomer:
		if err := r.parseOrderingCustomer(); err != nil {
			return err
		}
	case TagOrderingInstitution:
		if err := r.parseOrderingInstitution(); err != nil {
			return err
		}
	case TagIntermediaryInstitution:
		if err := r.parseIntermediaryInstitution(); err != nil {
			return err
		}
	case TagInstitutionAccount:
		if err := r.parseInstitutionAccount(); err != nil {
			return err
		}
	case TagBeneficiaryCustomer:
		if err := r.parseBeneficiaryCustomer(); err != nil {
			return err
		}
	case TagRemittance:
		if err := r.parseRemittance(); err != nil {
			return err
		}
	case TagSenderToReceiver:
		if err := r.parseSenderToReceiver(); err != nil {
			return err
		}
	case TagUnstructuredAddenda:
		if err := r.parseUnstructuredAddenda(); err != nil {
			return err
		}
	case TagRelatedRemittance:
		if err := r.parseRelatedRemittance(); err != nil {
			return err
		}
	case TagRemittanceOriginator:
		if err := r.parseRemittanceOriginator(); err != nil {
			return err
		}
	case TagRemittanceBeneficiary:
		if err := r.parseRemittanceBeneficiary(); err != nil {
			return err
		}
	case TagPrimaryRemittanceDocument:
		if err := r.parsePrimaryRemittanceDocument(); err != nil {
			return err
		}
	case TagActualAmountPaid:
		if err := r.parseActualAmountPaid(); err != nil {
			return err
		}
	case TagGrossAmountRemittanceDocument:
		if err := r.parseGrossAmountRemittanceDocument(); err != nil {
			return err
		}
	case TagAmountNegotiatedDiscount:
		if err := r.parseAmountNegotiatedDiscount(); err != nil {
			return err
		}
	case TagAdjustment:
		if err := r.parseAdjustment(); err != nil {
			return err
		}
	case TagDateRemittanceDocument:
		if err := r.parseDateRemittanceDocument(); err != nil {
			return err
		}
	case TagSecondaryRemittanceDocument:
		if err := r.parseSecondaryRemittanceDocument(); err != nil {
			return err
		}
	case TagRemittanceFreeText:
		if err := r.parseRemittanceFreeText(); err != nil {
			return err
		}
	case TagServiceMessage:
		if err := r.parseServiceMessage(); err != nil {
			return err
		}
	case TagReceiptTimeStamp:
		if err := r.parseReceiptTimeStamp(); err != nil {
			return err
		}
	default:
		return NewErrInvalidTag(r.line[:6])
	}
	return nil
}

func (r *Reader) parseSenderSupplied() error {
	r.tagName = "SenderSupplied"
	if len(r.line) != 18 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(18, len(r.line))))
		return r.errors
	}
	ss := new(SenderSupplied)
	ss.Parse(r.line)
	if err := ss.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetSenderSupplied(ss)
	return nil
}

func (r *Reader) parseTypeSubType() error {
	r.tagName = "TypeSubType"
	if len(r.line) != 10 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(10, len(r.line))))
		return r.errors
	}
	tst := new(TypeSubType)
	tst.Parse(r.line)
	if err := tst.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetTypeSubType(tst)
	return nil
}

func (r *Reader) parseInputMessageAccountabilityData() error {
	r.tagName = "InputMessageAccountabilityData"
	if len(r.line) != 28 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(22, len(r.line))))
		return r.errors
	}
	imad := new(InputMessageAccountabilityData)
	imad.Parse(r.line)
	if err := imad.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetInputMessageAccountabilityData(imad)
	return nil
}

func (r *Reader) parseAmount() error {
	r.tagName = "Amount"
	if len(r.line) != 18 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(18, len(r.line))))
		return r.errors
	}
	amt := new(Amount)
	amt.Parse(r.line)
	if err := amt.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetAmount(amt)
	return nil
}

func (r *Reader) parseSenderDepositoryInstitution() error {
	r.tagName = "SenderDepositoryInstitution"
	if len(r.line) < 15 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(15, len(r.line))))
		return r.errors
	}
	sdi := new(SenderDepositoryInstitution)
	sdi.Parse(r.line)
	if err := sdi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetSenderDepositoryInstitution(sdi)
	return nil
}

func (r *Reader) parseReceiverDepositoryInstitution() error {
	r.tagName = "ReceiverDepositoryInstitution"
	if len(r.line) < 15 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(15, len(r.line))))
		return r.errors
	}
	rdi := new(ReceiverDepositoryInstitution)
	rdi.Parse(r.line)
	if err := rdi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetReceiverDepositoryInstitution(rdi)
	return nil
}

func (r *Reader) parseBusinessFunctionCode() error {
	r.tagName = "BusinessFunctionCode"
	if len(r.line) != 12 {
		r.errors.Add(r.parseError(NewTagWrongLengthErr(12, len(r.line))))
		return r.errors
	}
	bfc := new(BusinessFunctionCode)
	bfc.Parse(r.line)
	if err := bfc.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBusinessFunctionCode(bfc)
	return nil
}

func (r *Reader) parseSenderReference() error {
	r.tagName = "SenderReference"
	sr := new(SenderReference)
	sr.Parse(r.line)
	if err := sr.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetSenderReference(sr)
	return nil
}

func (r *Reader) parsePreviousMessageIdentifier() error {
	r.tagName = "PreviousMessageIdentifier"
	pmi := new(PreviousMessageIdentifier)
	pmi.Parse(r.line)
	if err := pmi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetPreviousMessageIdentifier(pmi)
	return nil
}

func (r *Reader) parseLocalInstrument() error {
	r.tagName = "LocalInstrument"
	li := new(LocalInstrument)
	li.Parse(r.line)
	if err := li.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetLocalInstrument(li)
	return nil
}

func (r *Reader) parsePaymentNotification() error {
	r.tagName = "PaymentNotification"
	pn := new(PaymentNotification)
	pn.Parse(r.line)
	if err := pn.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetPaymentNotification(pn)
	return nil
}

func (r *Reader) parseCharges() error {
	r.tagName = "Charges"
	c := new(Charges)
	c.Parse(r.line)
	if err := c.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetCharges(c)
	return nil
}

func (r *Reader) parseInstructedAmount() error {
	r.tagName = "InstructedAmount"
	ia := new(InstructedAmount)
	ia.Parse(r.line)
	if err := ia.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetInstructedAmount(ia)
	return nil
}

func (r *Reader) parseExchangeRate() error {
	r.tagName = "ExchangeRate"
	eRate := new(ExchangeRate)
	eRate.Parse(r.line)
	if err := eRate.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetExchangeRate(eRate)
	return nil
}

func (r *Reader) parseBeneficiaryIntermediaryFI() error {
	r.tagName = "BeneficiaryIntermediaryFI"
	bifi := new(BeneficiaryIntermediaryFI)
	bifi.Parse(r.line)
	if err := bifi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBeneficiaryIntermediaryFI(bifi)
	return nil
}

func (r *Reader) parseBeneficiaryFI() error {
	r.tagName = "BeneficiaryFI"
	bfi := new(BeneficiaryFI)
	bfi.Parse(r.line)
	if err := bfi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBeneficiaryFI(bfi)
	return nil
}

func (r *Reader) parseBeneficiary() error {
	r.tagName = "Beneficiary"
	ben := new(Beneficiary)
	ben.Parse(r.line)
	if err := ben.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBeneficiary(ben)
	return nil
}

func (r *Reader) parseBeneficiaryReference() error {
	r.tagName = "BeneficiaryReference"
	br := new(BeneficiaryReference)
	br.Parse(r.line)
	if err := br.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBeneficiaryReference(br)
	return nil
}

func (r *Reader) parseAccountDebitedDrawdown() error {
	r.tagName = "AccountDebitedDrawdown"
	debitDD := new(AccountDebitedDrawdown)
	debitDD.Parse(r.line)
	if err := debitDD.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetAccountDebitedDrawdown(debitDD)
	return nil
}

func (r *Reader) parseOriginator() error {
	r.tagName = "Originator"
	o := new(Originator)
	o.Parse(r.line)
	if err := o.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOriginator(o)
	return nil
}

func (r *Reader) parseOriginatorOptionF() error {
	r.tagName = "OriginatorOptionF"
	oof := new(OriginatorOptionF)
	oof.Parse(r.line)
	if err := oof.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOriginatorOptionF(oof)
	return nil
}

func (r *Reader) parseOriginatorFI() error {
	r.tagName = "OriginatorFI"
	ofi := new(OriginatorFI)
	ofi.Parse(r.line)
	if err := ofi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOriginatorFI(ofi)
	return nil
}

func (r *Reader) parseInstructingFI() error {
	r.tagName = "InstructingFI"
	ifi := new(InstructingFI)
	ifi.Parse(r.line)
	if err := ifi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetInstructingFI(ifi)
	return nil
}

func (r *Reader) parseAccountCreditedDrawdown() error {
	r.tagName = "AccountCreditedDrawdown"
	creditDD := new(AccountCreditedDrawdown)
	creditDD.Parse(r.line)
	if err := creditDD.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetAccountCreditedDrawdown(creditDD)
	return nil
}

func (r *Reader) parseOriginatorToBeneficiary() error {
	r.tagName = "OriginatorToBeneficiary"
	ob := new(OriginatorToBeneficiary)
	ob.Parse(r.line)
	if err := ob.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOriginatorToBeneficiary(ob)
	return nil
}

func (r *Reader) parseFIReceiverFI() error {
	r.tagName = "FIReceiverFI"
	firfi := new(FIReceiverFI)
	firfi.Parse(r.line)
	if err := firfi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIReceiverFI(firfi)
	return nil
}

func (r *Reader) parseFIDrawdownDebitAccountAdvice() error {
	r.tagName = "FIDrawdownDebitAccountAdvice"
	debitDDAdvice := new(FIDrawdownDebitAccountAdvice)
	debitDDAdvice.Parse(r.line)
	if err := debitDDAdvice.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIDrawdownDebitAccountAdvice(debitDDAdvice)
	return nil
}

func (r *Reader) parseFIIntermediaryFI() error {
	r.tagName = "FIIntermediaryFI"
	fiifi := new(FIIntermediaryFI)
	fiifi.Parse(r.line)
	if err := fiifi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIIntermediaryFI(fiifi)
	return nil
}

func (r *Reader) parseFIIntermediaryFIAdvice() error {
	r.tagName = "FIIntermediaryFIAdvice"
	fiifia := new(FIIntermediaryFIAdvice)
	fiifia.Parse(r.line)
	if err := fiifia.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIIntermediaryFIAdvice(fiifia)
	return nil
}

func (r *Reader) parseFIBeneficiaryFI() error {
	r.tagName = "FIBeneficiaryFI"
	fibfi := new(FIBeneficiaryFI)
	fibfi.Parse(r.line)
	if err := fibfi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIBeneficiaryFI(fibfi)
	return nil
}

func (r *Reader) parseFIBeneficiaryFIAdvice() error {
	r.tagName = "FIBeneficiaryFIAdvice"
	fibfia := new(FIBeneficiaryFIAdvice)
	fibfia.Parse(r.line)
	if err := fibfia.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIBeneficiaryFIAdvice(fibfia)
	return nil
}

func (r *Reader) parseFIBeneficiary() error {
	r.tagName = "FIBeneficiary"
	fib := new(FIBeneficiary)
	fib.Parse(r.line)
	if err := fib.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIBeneficiary(fib)
	return nil
}

func (r *Reader) parseFIBeneficiaryAdvice() error {
	r.tagName = "FIBeneficiaryAdvice"
	fiba := new(FIBeneficiaryAdvice)
	fiba.Parse(r.line)
	if err := fiba.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIBeneficiaryAdvice(fiba)
	return nil
}

func (r *Reader) parseFIPaymentMethodToBeneficiary() error {
	r.tagName = "FIPaymentMethodToBeneficiary"
	pm := new(FIPaymentMethodToBeneficiary)
	pm.Parse(r.line)
	if err := pm.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIPaymentMethodToBeneficiary(pm)
	return nil
}

func (r *Reader) parseFIAdditionalFiToFi() error {
	r.tagName = "FIAdditionalFiToFi"
	fifi := new(FIAdditionalFIToFI)
	fifi.Parse(r.line)
	if err := fifi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetFIAdditionalFIToFI(fifi)
	return nil
}

func (r *Reader) parseCurrencyInstructedAmount() error {
	r.tagName = "CurrencyInstructedAmount"
	cia := new(CurrencyInstructedAmount)
	cia.Parse(r.line)
	if err := cia.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetCurrencyInstructedAmount(cia)
	return nil
}

func (r *Reader) parseOrderingCustomer() error {
	r.tagName = "OrderingCustomer"
	oc := new(OrderingCustomer)
	oc.Parse(r.line)
	if err := oc.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOrderingCustomer(oc)
	return nil
}

func (r *Reader) parseOrderingInstitution() error {
	r.tagName = "OrderingInstitution"
	oi := new(OrderingInstitution)
	oi.Parse(r.line)
	if err := oi.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOrderingInstitution(oi)
	return nil
}

func (r *Reader) parseIntermediaryInstitution() error {
	r.tagName = "IntermediaryInstitution"
	ii := new(IntermediaryInstitution)
	ii.Parse(r.line)
	if err := ii.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetIntermediaryInstitution(ii)
	return nil
}

func (r *Reader) parseInstitutionAccount() error {
	r.tagName = "InstitutionAccount"
	iAccount := new(InstitutionAccount)
	iAccount.Parse(r.line)
	if err := iAccount.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetInstitutionAccount(iAccount)
	return nil
}

func (r *Reader) parseBeneficiaryCustomer() error {
	r.tagName = "BeneficiaryCustomer"
	bc := new(BeneficiaryCustomer)
	bc.Parse(r.line)
	if err := bc.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetBeneficiaryCustomer(bc)
	return nil
}

func (r *Reader) parseRemittance() error {
	r.tagName = "Remittance"
	ri := new(Remittance)
	ri.Parse(r.line)
	if err := ri.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetRemittance(ri)
	return nil
}

func (r *Reader) parseSenderToReceiver() error {
	r.tagName = "SenderToReceiver"
	sr := new(SenderToReceiver)
	sr.Parse(r.line)
	if err := sr.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetSenderToReceiver(sr)
	return nil
}

func (r *Reader) parseUnstructuredAddenda() error {
	r.tagName = "UnstructuredAddenda"
	ua := new(UnstructuredAddenda)
	ua.Parse(r.line)
	if err := ua.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetUnstructuredAddenda(ua)
	return nil
}

func (r *Reader) parseRelatedRemittance() error {
	r.tagName = "RelatedRemittance"
	rr := new(RelatedRemittance)
	rr.Parse(r.line)
	if err := rr.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetRelatedRemittance(rr)
	return nil
}

func (r *Reader) parseRemittanceOriginator() error {
	r.tagName = "RemittanceOriginator"
	ro := new(RemittanceOriginator)
	ro.Parse(r.line)
	if err := ro.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetRemittanceOriginator(ro)
	return nil
}

func (r *Reader) parseRemittanceBeneficiary() error {
	r.tagName = "RemittanceBeneficiary"
	rb := new(RemittanceBeneficiary)
	rb.Parse(r.line)
	if err := rb.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetRemittanceBeneficiary(rb)
	return nil
}

func (r *Reader) parsePrimaryRemittanceDocument() error {
	r.tagName = "PrimaryRemittanceDocument"
	prd := new(PrimaryRemittanceDocument)
	prd.Parse(r.line)
	if err := prd.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetPrimaryRemittanceDocument(prd)
	return nil
}

func (r *Reader) parseActualAmountPaid() error {
	r.tagName = "ActualAmountPaid"
	aap := new(ActualAmountPaid)
	aap.Parse(r.line)
	if err := aap.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetActualAmountPaid(aap)
	return nil
}

func (r *Reader) parseGrossAmountRemittanceDocument() error {
	r.tagName = "GrossAmountRemittanceDocument"
	gard := new(GrossAmountRemittanceDocument)
	gard.Parse(r.line)
	if err := gard.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetGrossAmountRemittanceDocument(gard)
	return nil
}

func (r *Reader) parseAmountNegotiatedDiscount() error {
	r.tagName = "AmountNegotiatedDiscount"
	nd := new(AmountNegotiatedDiscount)
	nd.Parse(r.line)
	if err := nd.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetAmountNegotiatedDiscount(nd)
	return nil
}

func (r *Reader) parseAdjustment() error {
	r.tagName = "Adjustment"
	adj := new(Adjustment)
	adj.Parse(r.line)
	if err := adj.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetAdjustment(adj)
	return nil
}

func (r *Reader) parseDateRemittanceDocument() error {
	r.tagName = "DateRemittanceDocument"
	drd := new(DateRemittanceDocument)
	drd.Parse(r.line)
	if err := drd.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetDateRemittanceDocument(drd)
	return nil
}

func (r *Reader) parseSecondaryRemittanceDocument() error {
	r.tagName = "SecondaryRemittanceDocument"
	srd := new(SecondaryRemittanceDocument)
	srd.Parse(r.line)
	if err := srd.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetSecondaryRemittanceDocument(srd)
	return nil
}

func (r *Reader) parseRemittanceFreeText() error {
	r.tagName = "RemittanceFreeText"
	rft := new(RemittanceFreeText)
	rft.Parse(r.line)
	if err := rft.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetRemittanceFreeText(rft)
	return nil
}

func (r *Reader) parseServiceMessage() error {
	r.tagName = "ServiceMessage"
	sm := new(ServiceMessage)
	sm.Parse(r.line)
	if err := sm.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetServiceMessage(sm)
	return nil
}

func (r *Reader) parseMessageDisposition() error {
	r.tagName = "MessageDisposition"
	md := new(MessageDisposition)
	md.Parse(r.line)
	if err := md.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetMessageDisposition(md)
	return nil
}

func (r *Reader) parseReceiptTimeStamp() error {
	r.tagName = "ReceiptTimeStamp"
	rts := new(ReceiptTimeStamp)
	rts.Parse(r.line)
	if err := rts.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetReceiptTimeStamp(rts)
	return nil
}

func (r *Reader) parseOutputMessageAccountabilityData() error {
	r.tagName = "OutputMessageAccountabilityData"
	omad := new(OutputMessageAccountabilityData)
	omad.Parse(r.line)
	if err := omad.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetOutputMessageAccountabilityData(omad)
	return nil
}

func (r *Reader) parseErrorWire() error {
	r.tagName = "ErrorWire"
	ew := new(ErrorWire)
	ew.Parse(r.line)
	if err := ew.Validate(); err != nil {
		return r.parseError(err)
	}
	r.currentFEDWireMessage.SetErrorWire(ew)
	return nil
}
