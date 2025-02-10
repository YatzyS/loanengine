package repo

import (
	"context"

	"github.com/loanengine/internal/dao"
)

type LoanRepo interface {
	GetList(ctx context.Context, query *dao.LoanStateTable, page int, offset int) ([]*dao.LoanTable, error)
	GetState(ctx context.Context, query *dao.LoanTable) (string, error)
	ProposeLoan(ctx context.Context, query *dao.LoanTable) error
	ApproveLoan(ctx context.Context, query *dao.LoanStateTable) error
	DisburseLoan(ctx context.Context, query *dao.LoanStateTable) error
	Invest(ctx context.Context, query *dao.LoanInvestmentTable) (string, error)
	GetListOfInvestors(ctx context.Context, loanID string) ([]string, error)
	UpdateAgreementLink(ctx context.Context, loanId string, agreementLink string) error
}

type loanRepo struct{}

// GetListOfInvestors implements LoanRepo.
func (l *loanRepo) GetListOfInvestors(ctx context.Context, loanID string) ([]string, error) {
	panic("unimplemented")
}

// ApproveLoan implements LoanRepo.
func (l *loanRepo) ApproveLoan(ctx context.Context, query *dao.LoanStateTable) error {
	panic("unimplemented")
}

// DisburseLoan implements LoanRepo.
func (l *loanRepo) DisburseLoan(ctx context.Context, query *dao.LoanStateTable) error {
	panic("unimplemented")
}

// GetList implements LoanRepo.
func (l *loanRepo) GetList(ctx context.Context, query *dao.LoanStateTable, page int, offset int) ([]*dao.LoanTable, error) {
	panic("unimplemented")
}

// GetState implements LoanRepo.
func (l *loanRepo) GetState(ctx context.Context, query *dao.LoanTable) (string, error) {
	panic("unimplemented")
}

// Invest implements LoanRepo.
func (l *loanRepo) Invest(ctx context.Context, query *dao.LoanInvestmentTable) (string, error) {
	panic("unimplemented")
}

// ProposeLoan implements LoanRepo.
func (l *loanRepo) ProposeLoan(ctx context.Context, query *dao.LoanTable) error {
	panic("unimplemented")
}

// UpdateAgreementLink implements LoanRepo.
func (l *loanRepo) UpdateAgreementLink(ctx context.Context, loanId string, agreementLink string) error {
	panic("unimplemented")
}

func NewLoanRepo() LoanRepo {
	return &loanRepo{}
}
