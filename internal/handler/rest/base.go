package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/loanengine/internal/service"
)

type RestHandler interface {
	// Propose is called when a new loan is proposed.
	Propose(c *gin.Context)

	// Approve is called when the loan gets approved by field agent by sending the required details.
	Approve(c *gin.Context)

	// Invest is called when an investor wants to invest into a loan.
	Invest(c *gin.Context)

	// Disburse is called when the field agent uploads the proof to get the loan disbursed.
	Disburse(c *gin.Context)

	// GetState provides the current state of the loan.
	GetState(c *gin.Context)

	// GetList can be used to get a list of loans in any given state. State can be passed in the URL itself.
	// If no state is passed it will by default return list of all the loans.
	GetList(c *gin.Context)
}

type restHandler struct {
	loanService service.LoanService
}

func NewRestHandler(loanService service.LoanService) RestHandler {
	return &restHandler{loanService: loanService}
}
