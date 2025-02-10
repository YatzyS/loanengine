package constants

type LoanState string

const (
	PROPOSED LoanState = "PROPOSED"
	APPROVED LoanState = "APPROVED"
	INVESTED LoanState = "INVESTED"
	DISBURSED LoanState = "DISBURSED"
)