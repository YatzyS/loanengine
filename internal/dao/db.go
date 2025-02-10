package dao

import (
	"time"

	"github.com/loanengine/internal/common/constants"
)

type LoanTable struct {
	LoanID          string  `xorm:"text not null unique 'loan_id'"`
	BorrowerID      string  `xorm:"text not null 'borrower_id'"`
	PrincipleAmount int     `xorm:"int not null 'principle_amount'"`
	ROI             float64 `xorm:"decimal not null 'roi'"`  // For investor
	Rate            float64 `xorm:"decimal not null 'rate'"` // For borrower
	AgreementLink   string  `xorm:"text 'agreement_link'"`
	Version         int     `xorm:"version"`
}

func (l *LoanTable) TableName() string {
	return constants.LOAN_TABLE
}

type LoanStateTable struct {
	LoanID        string    `xorm:"text not null unique 'loan_id'"`
	LoanState     string    `xorm:"text not null 'loan_state'"`
	ApproveEmpID  string    `xorm:"text 'approve_emp_id'"`
	ApproveProof  string    `xorm:"text 'approve_proof'"`
	ApproveDate   time.Time `xorm:"datetime 'approve_time'"`
	DisburseEmpID string    `xorm:"text 'disburse_emp_id'"`
	DisburseProof string    `xorm:"text 'disburse_proof'"`
	DisburseDate  time.Time `xorm:"datetime 'disburse_time'"`
	Version       int       `xorm:"version"`
}

func (l *LoanStateTable) TableName() string {
	return constants.LOAN_STATE_TABLE
}

type LoanInvestmentTable struct {
	LoanID     string `xorm:"text not null unique 'loan_id'"`
	InvestorID string `xorm:"text not null 'investor_id'"`
	Amount     int    `xorm:"int not null 'amount'"`
	Version    int    `xorm:"version"`
}

func (l *LoanInvestmentTable) TableName() string {
	return constants.LOAN_INVESTMENT_TABLE
}

type LoanAmountTable struct {
	LoanID         string `xorm:"text not null unique 'loan_id'"`
	RequiredAmount int    `xorm:"int not null 'required_amount'"`
	InvestedAmount int    `xorm:"int 'invested_amount'"`
	Version        int    `xorm:"version"`
}

func (l *LoanAmountTable) TableName() string {
	return constants.LOAN_AMOUNT_TABLE
}
