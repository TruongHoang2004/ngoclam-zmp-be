package middleware

import (
	"fmt"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next() // process request

		// Check if any errors were added to context
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors {
				fmt.Println("Error occurred:", e.Err)
				if httpErr, ok := e.Err.(*common.HTTPException); ok {
					ctx.JSON(httpErr.Status, gin.H{"error": httpErr.Message})
					return
				}
			}
			// fallback if error not HTTPException
			ctx.JSON(500, gin.H{"error": ctx.Errors[0].Error()})
		}
	}
}
