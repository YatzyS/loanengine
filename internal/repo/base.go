package repo

import (
	"context"
	"fmt"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/loanengine/internal/common/constants"
	"github.com/loanengine/internal/dao"
	"github.com/sirupsen/logrus"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type LoanRepo interface {
	GetList(ctx context.Context, query *dao.LoanStateTable, limit int, offset int) ([]*dao.LoanTable, error)
	GetState(ctx context.Context, query *dao.LoanStateTable) (string, error)
	ProposeLoan(ctx context.Context, query *dao.LoanTable) error
	ApproveLoan(ctx context.Context, query *dao.LoanStateTable) error
	DisburseLoan(ctx context.Context, query *dao.LoanStateTable) error
	Invest(ctx context.Context, query *dao.LoanInvestmentTable) (string, error)
	GetListOfInvestors(ctx context.Context, loanID string) ([]string, error)
	UpdateAgreementLink(ctx context.Context, loanId string, agreementLink string) error
}

type loanRepo struct {
	engine *xorm.Engine
	mu     sync.RWMutex
}

// GetListOfInvestors implements LoanRepo.
func (l *loanRepo) GetListOfInvestors(ctx context.Context, loanID string) ([]string, error) {
	var rows []dao.LoanInvestmentTable
	err := l.engine.Context(ctx).Find(&rows, &dao.LoanInvestmentTable{LoanID: loanID})
	if err != nil {
		return nil, fmt.Errorf("get list of investors error: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("couldnt find any investor for the loan: %s", loanID)
	}
	resp := make([]string, len(rows))
	for idx, val := range rows {
		resp[idx] = val.InvestorID
	}
	return resp, nil
}

// ApproveLoan implements LoanRepo.
func (l *loanRepo) ApproveLoan(ctx context.Context, query *dao.LoanStateTable) error {
	loanStateRow := &dao.LoanStateTable{
		LoanID: query.LoanID,
	}
	ok, err := l.engine.Get(loanStateRow)
	if err != nil {
		return fmt.Errorf("get details for approve loan error: %w", err)
	}
	if !ok {
		return fmt.Errorf("loan details for approve loan not found")
	}
	if loanStateRow.LoanState != string(constants.PROPOSED) {
		return fmt.Errorf("invalid state trasnsition to approve")
	}
	query.Version = loanStateRow.Version
	_, err = l.engine.Context(ctx).Where("loan_id=?",query.LoanID).Update(query)
	if err != nil {
		return fmt.Errorf("approve loan error: %w", err)
	}
	return nil
}

// DisburseLoan implements LoanRepo.
func (l *loanRepo) DisburseLoan(ctx context.Context, query *dao.LoanStateTable) error {
	loanStateRow := &dao.LoanStateTable{
		LoanID: query.LoanID,
	}
	ok, err := l.engine.Get(loanStateRow)
	if err != nil {
		return fmt.Errorf("get details for disburse loan error: %w", err)
	}
	if !ok {
		return fmt.Errorf("loan details for disburse loan not found")
	}
	if loanStateRow.LoanState != string(constants.INVESTED) {
		return fmt.Errorf("invalid state trasnsition to disburse")
	}
	query.Version = loanStateRow.Version
	_, err = l.engine.Context(ctx).Where("loan_id=?",query.LoanID).Update(query)
	if err != nil {
		return fmt.Errorf("disburse loan error: %w", err)
	}
	return nil
}

// GetList implements LoanRepo.
func (l *loanRepo) GetList(ctx context.Context, query *dao.LoanStateTable, limit int, offset int) ([]*dao.LoanTable, error) {
	/*
		Get list of loan Ids with required state
	*/
	var list []*dao.LoanTable
	qb :=  builder.Select("loan_id").From(constants.LOAN_STATE_TABLE)
	if query.LoanState != "" {
		qb = qb.Where(builder.Eq{"loan_state": query.LoanState})
	}
	err := l.engine.Context(ctx).Limit(limit, offset).In("loan_id", qb).Find(&list)
	if err != nil {
		return nil, fmt.Errorf("get list error: %w", err)
	}
	return list, nil
}

// GetState implements LoanRepo.
func (l *loanRepo) GetState(ctx context.Context, query *dao.LoanStateTable) (string, error) {
	ok, err := l.engine.Context(ctx).Get(query)
	if err != nil {
		return "", fmt.Errorf("get state error: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("loan not found with id: %s", query.LoanID)
	}
	return query.LoanState, nil
}

// Invest implements LoanRepo.
func (l *loanRepo) Invest(ctx context.Context, query *dao.LoanInvestmentTable) (string, error) {
	/*
		Read lock on the code to get the invested amount
		Write lock while updating the inevestment.
	*/
	session := l.engine.NewSession()
	defer session.Close()
	l.mu.Lock()
	defer l.mu.Unlock()
	lastState := &dao.LoanStateTable{LoanID: query.LoanID}
	ok, err := l.engine.Context(ctx).Get(lastState)
	if err != nil {
		session.Rollback()
		return "", fmt.Errorf("get loan state error:%w", err)
	}
	if !ok {
		session.Rollback()
		return "", fmt.Errorf("loan state not found for loan id")
	}
	if lastState.LoanState == string(constants.INVESTED) {
		session.Rollback()
		return "", fmt.Errorf("loan investment is complete")
	}
	if lastState.LoanState != string(constants.APPROVED) {
		session.Rollback()
		return "", fmt.Errorf("loan not in state for investment")
	}
	row := &dao.LoanAmountTable{LoanID: query.LoanID}
	ok, err = l.engine.Context(ctx).Get(row)
	if err != nil {
		session.Rollback()
		return "", fmt.Errorf("get loan amount error:%w", err)
	}
	if !ok {
		session.Rollback()
		return "", fmt.Errorf("loan amount not found for loan id")
	}
	if row.InvestedAmount+query.Amount > row.RequiredAmount {
		session.Rollback()
		return "", fmt.Errorf("investment amount higher then the remaining amount: %d", row.RequiredAmount-row.InvestedAmount)
	}
	_, err = l.engine.Context(ctx).Insert(query)
	if err != nil {
		session.Rollback()
		return "", fmt.Errorf("insert loan investment error: %w", err)
	}
	q := &dao.LoanAmountTable{InvestedAmount: row.InvestedAmount + query.Amount, Version: row.Version}
	_, err = l.engine.Context(ctx).Where("loan_id=?",row.LoanID).Update(q)
	if err != nil {
		session.Rollback()
		return "", fmt.Errorf("update loan amount error: %w", err)
	}
	stateRow := &dao.LoanStateTable{LoanID: query.LoanID}
	ok, err = l.engine.Context(ctx).Get(stateRow)
	if err != nil {
		session.Rollback()
		return "", fmt.Errorf("get loan state error:%w", err)
	}
	if !ok {
		session.Rollback()
		return "", fmt.Errorf("loan state not found for loan id")
	}
	if q.InvestedAmount == row.RequiredAmount {
		row := &dao.LoanStateTable{LoanState: string(constants.INVESTED), Version: stateRow.Version}
		_, err = l.engine.Context(ctx).Where("loan_id=?",stateRow.LoanID).Update(row)
		if err != nil {
			session.Rollback()
			return "", fmt.Errorf("update loan state error: %w", err)
		}
	}
	err = session.Commit()
	if err != nil {
		return "", fmt.Errorf("invest trasnaction commit error: %w", err)
	}
	return stateRow.LoanState, nil
}

// ProposeLoan implements LoanRepo.
func (l *loanRepo) ProposeLoan(ctx context.Context, query *dao.LoanTable) error {
	session := l.engine.NewSession()
	defer session.Close()
	_, err := l.engine.Context(ctx).Insert(query)
	if err != nil {
		session.Rollback()
		return fmt.Errorf("insert loan error: %w", err)
	}
	q1 := &dao.LoanStateTable{
		LoanID:    query.LoanID,
		LoanState: string(constants.PROPOSED),
	}
	_, err = l.engine.Context(ctx).Insert(q1)
	if err != nil {
		session.Rollback()
		return fmt.Errorf("insert loan status error: %w", err)
	}
	q2 := &dao.LoanAmountTable{
		LoanID:         query.LoanID,
		RequiredAmount: query.PrincipleAmount,
	}
	_, err = l.engine.Context(ctx).Insert(q2)
	if err != nil {
		session.Rollback()
		return fmt.Errorf("insert loan amount error: %w", err)
	}
	err = session.Commit()
	if err != nil {
		return fmt.Errorf("propose trasnaction commit error: %w", err)
	}
	return nil
}

// UpdateAgreementLink implements LoanRepo.
func (l *loanRepo) UpdateAgreementLink(ctx context.Context, loanId string, agreementLink string) error {
	loanRow := &dao.LoanTable{
		LoanID: loanId,
	}
	ok, err := l.engine.Context(ctx).Get(loanRow)
	if err != nil {
		return fmt.Errorf("get details for agreement error: %w", err)
	}
	if !ok {
		return fmt.Errorf("loan details for agreement not found")
	}
	q := &dao.LoanTable{AgreementLink: agreementLink, Version: loanRow.Version}
	_, err = l.engine.Context(ctx).Where("loan_id=?",loanRow.LoanID).Update(q)
	if err != nil {
		return fmt.Errorf("update agreement error: %w", err)
	}
	return nil
}

func NewLoanRepo() LoanRepo {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		constants.DB_HOST, constants.DB_PORT, constants.DB_USER, constants.DB_PASS, constants.DB_NAME)
	var err error
	en, err := xorm.NewEngine("pgx", dbinfo)
	if err != nil {
		log.Fatalf("engine creation failed", err)
	}
	en.ShowSQL(true)
	err = en.Ping()
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("Successfully connected")
	return &loanRepo{engine: en, mu: sync.RWMutex{}}
}
