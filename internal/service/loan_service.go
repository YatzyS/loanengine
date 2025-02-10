package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/loanengine/internal/common/constants"
	"github.com/loanengine/internal/dao"
	"github.com/loanengine/internal/repo"
)

type LoanService interface {
	Propose(ctx context.Context, loan *dao.Loan) error
	Approve(ctx context.Context, approve *dao.ApproveDetails) error
	Invest(ctx context.Context, invest *dao.LoanInvest) error
	Disburse(ctx context.Context, disburse *dao.ApproveDetails) error
	GetState(ctx context.Context, loanId string) (*dao.LoanStateResponse, error)
	GetList(ctx context.Context, page int, offset int, state string) (*dao.GetListResponse, error)
}

type loanService struct {
	repo                repo.LoanRepo
	agreementGenerator  AgreementGenerator
	notificationService NotificationService
}

func (l *loanService) Approve(ctx context.Context, approve *dao.ApproveDetails) error {
	query := &dao.LoanStateTable{
		LoanID:       approve.LoanId,
		ApproveEmpID: approve.EmpId,
		ApproveProof: approve.ImagePath,
		ApproveDate:  approve.Date,
	}
	err := l.repo.ApproveLoan(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (l *loanService) GetList(ctx context.Context, page int, offset int, state string) (*dao.GetListResponse, error) {
	query := &dao.LoanStateTable{}
	if state != "" {
		query.LoanState = state
	}
	dbResp, err := l.repo.GetList(ctx, query, page, offset)
	if err != nil {
		return nil, err
	}
	resp := &dao.GetListResponse{
		Loans: make([]*dao.LoanListDetails, len(dbResp)),
	}
	for idx, val := range dbResp {
		resp.Loans[idx] = &dao.LoanListDetails{
			LoanId:          val.LoanID,
			PrincipleAmount: val.PrincipleAmount,
			Rate:            val.Rate,
			ROI:             val.ROI,
		}
	}
	return resp, nil
}

func (l *loanService) GetState(ctx context.Context, loanId string) (*dao.LoanStateResponse, error) {
	query := &dao.LoanTable{LoanID: loanId}
	// TODO: Check if facade is required as it queries on two tables.
	state, err := l.repo.GetState(ctx, query)
	if err != nil {
		return nil, err
	}
	resp := &dao.LoanStateResponse{
		LoanId:    loanId,
		LoanState: state,
	}
	return resp, nil
}

func (l *loanService) Invest(ctx context.Context, invest *dao.LoanInvest) error {
	/*
		Check if the invested amount is already met.
		If so return with an error.
		Use mutext to make sure the operation is done right.
		Update the status to invested if the investment amount is matched with principle amount.
	*/
	query := &dao.LoanInvestmentTable{
		LoanID:     invest.LoanId,
		InvestorID: invest.InvestorId,
		Amount:     invest.Amount,
	}
	state, err := l.repo.Invest(ctx, query)
	if err != nil {
		return err
	}
	if constants.LoanState(state) != constants.INVESTED {
		return nil
	}
	err = l.agreementGenerator.Generate(ctx, query.LoanID)
	if err != nil {
		return err
	}
	err = l.notificationService.NotifyForAgreement(ctx, query.LoanID)
	if err != nil {
		return err
	}
	return nil
}

func (l *loanService) Disburse(ctx context.Context, disburse *dao.ApproveDetails) error {
	query := &dao.LoanStateTable{
		LoanID:        disburse.LoanId,
		DisburseEmpID: disburse.EmpId,
		DisburseProof: disburse.ImagePath,
		DisburseDate:  disburse.Date,
	}
	err := l.repo.DisburseLoan(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (l *loanService) Propose(ctx context.Context, loan *dao.Loan) error {
	query := &dao.LoanTable{
		LoanID:          uuid.NewString(),
		BorrowerID:      loan.BorrowerId,
		PrincipleAmount: loan.PrincipleAmount,
		ROI:             loan.ROI,
		Rate:            loan.Rate,
	}
	// This will do two operations in transaction.
	// 1. Add new loan to Loan table.
	// 2. Add new loan in LoanState table, with initial state as PROPOSED.
	err := l.repo.ProposeLoan(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func NewLoanService(repo repo.LoanRepo, agreementGenerator AgreementGenerator, notificationService NotificationService) LoanService {
	return &loanService{repo: repo, agreementGenerator: agreementGenerator, notificationService: notificationService}
}
