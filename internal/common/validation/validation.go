package validation

import (
	"fmt"

	"github.com/loanengine/internal/dao"
	"go.uber.org/multierr"
)

func CheckLoanDetails(req *dao.Loan) error {
	if req.BorrowerId == "" {
		return fmt.Errorf("validation: invalid borrower_id")
	}
	// TODO: Check if we need validation on ROI, Rate and Principle Amount
	return nil
}

func CheckInvestRequest(req *dao.LoanInvest) error {
	var err error
	if req.InvestorId == "" {
		err = multierr.Append(err, fmt.Errorf("validation: invalid investor_id"))
	}
	if req.LoanId == "" {
		err = multierr.Append(err, fmt.Errorf("validation: invalid loan_id"))
	}
	// TODO: Check if we need validation on Amount
	return err
}

func CheckApproveReq(req *dao.ApproveDetails) error {
	var err error
	if req.EmpId == "" {
		err = multierr.Append(err, fmt.Errorf("validation: invalid emp_id"))
	}
	if req.LoanId == "" {
		err = multierr.Append(err, fmt.Errorf("validation: invalid loan_id"))
	}
	return err
}
