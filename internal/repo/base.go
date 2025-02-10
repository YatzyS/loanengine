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
	UpdateAgreementLink(ctx context.Context, loanId string, agreementLink string) error
}
