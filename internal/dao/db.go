package dao

import "time"

type LoanTable struct {
	LoanID          string
	BorrowerID      string
	PrincipleAmount int
	ROI             float64 // For investor
	Rate            float64 // For borrower
	AgreementLink   string
}

type LoanStateTable struct {
	LoanID        string
	LoanState     string
	ApproveEmpID  string
	ApproveProof  string
	ApproveDate   time.Time
	DisburseEmpID string
	DisburseProof string
	DisburseDate  time.Time
}

type LoanInvestmentTable struct {
	LoanID     string
	InvestorID string
	Amount     string
}

type LoanAmountTable struct {
	LoanID         string
	RequiredAmount int
	InvestedAmount int
}
