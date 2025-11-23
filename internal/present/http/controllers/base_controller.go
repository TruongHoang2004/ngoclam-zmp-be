package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type baseController struct {
	validate *validator.Validate
}

func NewBaseController(validate *validator.Validate) *baseController {
	return &baseController{
		validate: validate,
	}
}

func (b *baseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (b *baseController) ErrorData(c *gin.Context, err *common.Error) {
	log.IErr(c.Request.Context(), err)
	c.JSON(err.GetHttpStatus(), common.ConvertErrorToResponse(err))
}

func (b *baseController) BindAndValidateRequest(c *gin.Context, req interface{}) *common.Error {
	if err := c.BindUri(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return common.ErrBadRequest(c).SetDetail(err.Error()).SetSource(common.CurrentService)
	}
	if err := c.Bind(req); err != nil {
		log.Warn(c, "bind request err, err:[%s]", err)
		return common.ErrBadRequest(c).SetDetail(err.Error()).SetSource(common.CurrentService)
	}
	return b.ValidateRequest(c, req)
}

func (b *baseController) ValidateRequest(ctx context.Context, req interface{}) *common.Error {
	err := b.validate.Struct(req)

	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Error(ctx, "Cannot parse validate error: %+v", err)
			return common.ErrSystemError(ctx, "ValidateFailed").SetDetail(err.Error()).SetSource(common.CurrentService)
		}
		var filedErrors []string
		for _, errValidate := range errs {
			log.Debug(ctx, "field invalid, err:[%s]", errValidate.Field())
			filedErrors = append(filedErrors, errValidate.Error())
		}
		str := strings.Join(filedErrors, ",")
		log.Warn(ctx, "invalid request, err:[%s]", err.Error())
		return common.ErrBadRequest(ctx).SetDetail(fmt.Sprintf("field invalidate [%s]", str)).SetSource(common.CurrentService)
	}
	return nil
}
