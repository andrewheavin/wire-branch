package wire

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// mockCharges creates a Charges
func mockCharges() *Charges {
	c := NewCharges()
	c.ChargeDetails = "B"
	c.SendersChargesOne = "USD0,99"
	c.SendersChargesTwo = "USD2,99"
	c.SendersChargesThree = "USD3,99"
	c.SendersChargesFour = "USD1,00"
	return c
}

// TestMockCharges validates mockCharges
func TestMockCharges(t *testing.T) {
	c := mockCharges()

	require.NoError(t, c.Validate(), "mockCharges does not validate and will break other tests")
}

// TestChargeDetailsValid validates ChargeDetails is valid
func TestPaymentNotificationIndicatorValid(t *testing.T) {
	c := mockCharges()
	c.ChargeDetails = "F"

	err := c.Validate()

	require.EqualError(t, err, fieldError("ChargeDetails", ErrChargeDetails, c.ChargeDetails).Error())
}

func TestChargesCrash(t *testing.T) {
	c := &Charges{}
	c.Parse("{3700}") // invalid, caused a fuzz crash

	require.Empty(t, c.tag)
	require.Empty(t, c.ChargeDetails)
}

// TestStringChargesVariableLength parses using variable length
func TestStringChargesVariableLength(t *testing.T) {
	var line = "{3700}"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseCharges()
	expected := r.parseError(NewTagMinLengthErr(7, len(r.line))).Error()
	require.EqualError(t, err, expected)

	line = "{3700}B                                                            NNN"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseCharges()
	require.EqualError(t, err, r.parseError(NewTagMaxLengthErr()).Error())

	line = "{3700}B*****"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseCharges()
	require.EqualError(t, err, r.parseError(NewTagMaxLengthErr()).Error())

	line = "{3700}B*"
	r = NewReader(strings.NewReader(line))
	r.line = line

	err = r.parseCharges()
	require.Equal(t, err, nil)
}

// TestStringChargesOptions validates string() with options
func TestStringChargesOptions(t *testing.T) {
	var line = "{3700}B*"
	r := NewReader(strings.NewReader(line))
	r.line = line

	err := r.parseCharges()
	require.Equal(t, err, nil)

	str := r.currentFEDWireMessage.Charges.String()
	require.Equal(t, str, "{3700}B                                                            ")

	str = r.currentFEDWireMessage.Charges.String(true)
	require.Equal(t, str, "{3700}B*")
}
