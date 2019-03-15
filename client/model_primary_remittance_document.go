/*
 * WIRE API
 *
 * Moov WIRE () implements an HTTP API for creating, parsing and validating WIRE files.
 *
 * API version: v1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type PrimaryRemittanceDocument struct {
	// DocumentTypeCode  * `AROI` - Accounts Receivable Open Item * `BOLD` - Bill of Lading Shipping Notice * `CINV` - Commercial Invoice * `CMCN` - Commercial Contract * `CNFA` - Credit Note Related to Financial Adjustment * `CREN` - Credit Note * `DEBN` - Debit Note * `DISP` - Dispatch Advice * `DNFA` - Debit Note Related to Financial Adjustment HIRI Hire Invoice * `MSIN` - Metered Service Invoice * `PROP` - Proprietary Document Type * `PUOR` - Purchase Order * `SBIN` - Self Billed Invoice * `SOAC` - Statement of Account * `TSUT` - Trade Services Utility Transaction VCHR Voucher 
	DocumentTypeCode string `json:"documentTypeCode,omitempty"`
	// ProprietaryDocumentTypeCode
	ProprietaryDocumentTypeCode string `json:"proprietaryDocumentTypeCode,omitempty"`
	// DocumentIdentificationNumber
	DocumentIdentificationNumber string `json:"documentIdentificationNumber,omitempty"`
	// Issuer
	Issuer string `json:"issuer,omitempty"`
}
