package common

import (
	"context"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/constant"
)

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ""
	if ctx.Value(constant.TraceIdName) != nil {
		traceId = ctx.Value(constant.TraceIdName).(string)
	}
	return traceId
}
