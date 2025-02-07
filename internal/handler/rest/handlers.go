package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/loanengine/internal/common/constants"
	"github.com/loanengine/internal/common/validation"
	"github.com/loanengine/internal/dao"
	log "github.com/sirupsen/logrus"
)

func (r restHandler) Propose(c *gin.Context) {
	/*
		1. Get Data into struct
		2. Validate Data
		3. Call Service layer
	*/
	req := &dao.Loan{}
	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := validation.CheckLoanDetails(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Propose(c, req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r restHandler) Approve(c *gin.Context) {
	/*
		1. Save the uploaded file to a local folder.
		2. Name it using UUID.
		3. Save the path in a string and pass it on to service for storing.
	*/
	req := &dao.ApproveDetails{}
	if err := c.ShouldBind(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dstPath := constants.APPROVE_FILE_UPLOAD_PATH + req.LoanId + "/" + uuid.New().String() + ".jpg"
	err = c.SaveUploadedFile(photo, dstPath)
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.ImagePath = dstPath
	if err := validation.CheckApproveReq(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Approve(c, req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r restHandler) Invest(c *gin.Context) {
	req := &dao.LoanInvest{}
	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := validation.CheckInvestRequest(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Invest(c, req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r restHandler) Disburse(c *gin.Context) {
	/*
		1. Save the uploaded file to a local folder.
		2. Name it using UUID.
		3. Save the path in a string and pass it on to service for storing.
	*/
	// TODO: Try to come up with a method as there is repeating code
	req := &dao.ApproveDetails{}
	if err := c.ShouldBind(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	dstPath := constants.DISBURSE_FILE_UPLOAD_PATH + req.LoanId + "/" + uuid.New().String() + ".jpg"
	err = c.SaveUploadedFile(photo, dstPath)
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.ImagePath = dstPath
	if err := validation.CheckApproveReq(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Approve(c, req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r restHandler) GetState(c *gin.Context) {
	loanId := c.Param("id")
	if loanId == "" {
		err := fmt.Errorf("loan engine: loan id cannot be empty")
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.GetState(c, loanId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (r restHandler) GetList(c *gin.Context) {
	state := c.Param("state")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.GetList(c, page, offset, state); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
