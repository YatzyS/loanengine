package dao

import "time"

type Loan struct {
	BorrowerId      string  `json:"borrower_id"`
	PrincipleAmount int     `json:"principle_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
	AgreementLink   string  `json:"aggreement_link,omitempty"`
}

type LoanInvest struct {
	InvestorId string `json:"investor_id"`
	LoanId     string `json:"loan_id"`
	Amount     int    `json:"amount"`
}

type VerifyDetails struct {
	LoanId    string    `form:"loan_id"`
	EmpId     string    `form:"emp_id"`
	Date      time.Time `form:"date"`
	ImagePath string
}

type GetListResponse struct {
	Loans []*LoanListDetails
}

type LoanListDetails struct {
	LoanId          string  `json:"loan_id"`
	PrincipleAmount int     `json:"principle_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
}

type LoanStateResponse struct {
	LoanId    string `json:"loan_id"`
	LoanState string `json:"loan_state"`
}

type GenericResponse struct {
	Message string `json:"message"`
}
