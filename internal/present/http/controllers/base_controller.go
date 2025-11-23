package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils/casting"
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

func (b *baseController) GetUintParam(c *gin.Context, key string) (uint, *common.Error) {
	param := c.Param(key)
	if param == "" {
		log.Warn(c, "param %s is empty", key)
		return 0, common.ErrBadRequest(c).SetDetail(fmt.Sprintf("param %s is empty", key)).SetSource(common.CurrentService)
	}

	id, err := casting.StringToUint(param)
	if err != nil {
		log.Warn(c, "invalid param %s, err:[%s]", key, err)
		return 0, common.ErrBadRequest(c).SetDetail(fmt.Sprintf("invalid param %s", key)).SetSource(common.CurrentService)
	}
	return id, nil
}

func (b *baseController) GetFile(c *gin.Context, key string) (*multipart.FileHeader, *common.Error) {
	file, err := c.FormFile(key)
	if err != nil {
		log.Warn(c, "get file %s err, err:[%s]", key, err)
		return nil, common.ErrBadRequest(c).SetDetail(fmt.Sprintf("file %s is required", key)).SetSource(common.CurrentService)
	}
	return file, nil
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
