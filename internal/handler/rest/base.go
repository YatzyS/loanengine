package rest

import "github.com/gin-gonic/gin"

type RestHandler interface {
	Propose(c *gin.Context)
	Approve(c *gin.Context)
	Invest(c *gin.Context)
	Disburse(c *gin.Context)
	GetState(c *gin.Context)
	GetProposedLoans(c *gin.Context)
}

type restHandler struct{}

func (r restHandler) Propose(c *gin.Context) {
	panic("implement me")
}

func (r restHandler) Approve(c *gin.Context) {
	panic("implement me")
}

func (r restHandler) Invest(c *gin.Context) {
	panic("implement me")
}

func (r restHandler) Disburse(c *gin.Context) {
	panic("implement me")
}

func (r restHandler) GetState(c *gin.Context) {
	panic("implement me")
}

func (r restHandler) GetProposedLoans(c *gin.Context) {
	panic("implement me")
}

func NewRestHandler() RestHandler {
	return &restHandler{}
}
