package service

import (
	"context"

	"github.com/loanengine/internal/repo"
)

type NotificationService interface {
	// NotifyForAgreement queries the DB to get details about the investors and will notify them about the generated agreement and share the link.
	NotifyForAgreement(ctx context.Context, loanID string) error
}

type notificationService struct {
	repo *repo.LoanRepo
}

// NotifyForAgreement implements NotificationService.
func (n *notificationService) NotifyForAgreement(ctx context.Context, loanID string) error {
	return nil
}

func NewNotificationService() NotificationService {
	return &notificationService{}
}
