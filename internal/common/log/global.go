package log

import (
	"context"
	"fmt"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/helpers"
)

var globalLogger *logger

func Info(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(addCtxValue(ctx, msg), args...)
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(addCtxValue(ctx, msg), args...)
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(addCtxValue(ctx, msg), args...)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(addCtxValue(ctx, msg), args...)
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(msg, args...)
}

func IErr(ctx context.Context, err *common.Error) {
	if helpers.IsInternalError(err) {
		globalLogger.Error(addCtxValue(ctx, err.GetDetail()))
	} else if helpers.IsClientError(err) {
		globalLogger.Warn(addCtxValue(ctx, err.ToJSon()))
	}

}

func GetLogger() *logger {
	return globalLogger
}

func addCtxValue(ctx context.Context, msg string) string {
	traceId := common.GetTraceId(ctx)

	return fmt.Sprintf("%s, trace_id:[%s]", msg, traceId)
}
