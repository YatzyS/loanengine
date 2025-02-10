package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	resp := &dao.GenericResponse{}
	req := &dao.Loan{}
	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := validation.CheckLoanDetails(req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Propose(c, req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	resp.Message = constants.SUCCESS_MESSAGE
	c.JSON(http.StatusOK, resp)
}

func (r restHandler) Approve(c *gin.Context) {
	/*
		1. Save the uploaded file to a local folder.
		2. Name it using UUID.
		3. Save the path in a string and pass it on to service for storing.
	*/
	req := &dao.VerifyDetails{}
	resp := &dao.GenericResponse{}
	if err := c.ShouldBind(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	dstPath := constants.APPROVE_FILE_UPLOAD_PATH + req.LoanId + "/" + uuid.New().String() + ".jpg"
	err = c.SaveUploadedFile(photo, dstPath)
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	req.ImagePath = dstPath
	if err := validation.CheckApproveReq(req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if err := r.loanService.Approve(c, req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Message = constants.SUCCESS_MESSAGE
	c.JSON(http.StatusOK, resp)
}

func (r restHandler) Invest(c *gin.Context) {
	req := &dao.LoanInvest{}
	resp := &dao.GenericResponse{}
	if err := c.BindJSON(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := validation.CheckInvestRequest(req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := r.loanService.Invest(c, req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Message = constants.SUCCESS_MESSAGE
	c.JSON(http.StatusOK, resp)
}

func (r restHandler) Disburse(c *gin.Context) {
	/*
		1. Save the uploaded file to a local folder.
		2. Name it using UUID.
		3. Save the path in a string and pass it on to service for storing.
	*/
	// TODO: Try to come up with a method as there is repeating code
	req := &dao.VerifyDetails{}
	resp := &dao.GenericResponse{}
	if err := c.ShouldBind(req); err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	dstPath := constants.DISBURSE_FILE_UPLOAD_PATH + req.LoanId + "/" + uuid.New().String() + ".jpg"
	err = c.SaveUploadedFile(photo, dstPath)
	if err != nil {
		err = fmt.Errorf("loan engine:%w", err)
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	req.ImagePath = dstPath
	if err := validation.CheckApproveReq(req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if err := r.loanService.Disburse(c, req); err != nil {
		log.WithContext(c).Error(err)
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	resp.Message = constants.SUCCESS_MESSAGE
	c.JSON(http.StatusOK, resp)
}

func (r restHandler) GetState(c *gin.Context) {
	errResp := &dao.GenericResponse{}
	loanId := c.Param("id")
	if loanId == "" {
		err := fmt.Errorf("loan engine: loan id cannot be empty")
		log.WithContext(c).Error(err)
		errResp.Message = err.Error()
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	var err error
	resp := &dao.LoanStateResponse{}
	if resp, err = r.loanService.GetState(c, loanId); err != nil {
		log.WithContext(c).Error(err)
		errResp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (r restHandler) GetList(c *gin.Context) {
	errResp := &dao.GenericResponse{}
	state := c.Param("state")
	state = strings.Trim(state, "/")
	state = strings.ToUpper(state)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		err = fmt.Errorf("loan engine invalid limit:%w", err)
		log.WithContext(c).Error(err)
		errResp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		err = fmt.Errorf("loan engine invalid offset:%w", err)
		log.WithContext(c).Error(err)
		errResp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	resp := &dao.GetListResponse{}
	if resp, err = r.loanService.GetList(c, limit, offset, state); err != nil {
		log.WithContext(c).Error(err)
		errResp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
