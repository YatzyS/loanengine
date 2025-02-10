package service

import (
	"context"

	"github.com/loanengine/internal/dao"
)

type LoanService interface {
	Propose(ctx context.Context, loan *dao.Loan) error
	Approve(ctx context.Context, approve *dao.ApproveDetails) error
	Invest(ctx context.Context, invest *dao.LoanInvest) error
	GetState(ctx context.Context, loanId string) (*dao.LoanStateResponse, error)
	GetList(ctx context.Context, page int, offset int, state string) (*dao.GetListResponse, error)
}

type loanService struct{}

func (l *loanService) Approve(ctx context.Context, approve *dao.ApproveDetails) error {
	panic("unimplemented")
}

func (l *loanService) GetList(ctx context.Context, page int, offset int, state string) (*dao.GetListResponse, error) {
	panic("unimplemented")
}

func (l *loanService) GetState(ctx context.Context, loanId string) (*dao.LoanStateResponse, error) {
	panic("unimplemented")
}

func (l *loanService) Invest(ctx context.Context, invest *dao.LoanInvest) error {
	panic("unimplemented")
}

func (l *loanService) Propose(ctx context.Context, loan *dao.Loan) error {
	panic("unimplemented")
}

func NewLoanService() LoanService {
	return &loanService{}
}
