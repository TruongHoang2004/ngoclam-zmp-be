package common

import (
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/utils"
	"github.com/gin-gonic/gin"
)

func ParseUintParam(ctx *gin.Context, param string) (uint, error) {
	valueStr := ctx.Param(param)
	valueUint, err := utils.StringToUint(valueStr)
	if err != nil {
		return 0, err
	}
	return valueUint, nil
}
