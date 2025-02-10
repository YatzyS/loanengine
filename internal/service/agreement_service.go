package service

import (
	"context"

	"github.com/loanengine/internal/repo"
)

type AgreementGenerator interface {
	Generate(ctx context.Context, loanID string) error
}

type agreementGenerator struct {
	repo repo.LoanRepo
}

// Generate implements AgreementGenerator.
func (a *agreementGenerator) Generate(ctx context.Context, loanID string) error {
	/*
		Generate an agreement and update the loan with the agreement link.
	*/
	agreementLink := "test-link/test.pdf"
	err := a.repo.UpdateAgreementLink(ctx, loanID, agreementLink)
	if err != nil {
		return err
	}
	return nil
}

func NewAgreementGenerator() AgreementGenerator {
	return &agreementGenerator{}
}
