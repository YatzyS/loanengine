package service

import (
	"context"

	"github.com/loanengine/internal/repo"
	"github.com/sirupsen/logrus"
)

type NotificationService interface {
	// NotifyForAgreement queries the DB to get details about the investors and will notify them about the generated agreement and share the link.
	NotifyForAgreement(ctx context.Context, loanID string) error
}

type emailNotification struct {
	repo repo.LoanRepo
}

// NotifyForAgreement implements NotificationService.
func (e *emailNotification) NotifyForAgreement(ctx context.Context, loanID string) error {
	investorList, err := e.repo.GetListOfInvestors(ctx, loanID)
	if err != nil {
		return err
	}
	// Printing the investors list for now. This will be used to query investor details and notify them as required.
	logrus.WithContext(ctx).Infof("list of investors: ", investorList)
	return nil
}

func NewNotificationService(repo repo.LoanRepo) NotificationService {
	return &emailNotification{repo: repo}
}
