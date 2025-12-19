package middleware

import (
	"net/http"
	"strings"

	"github.com/TruongHoang2004/ngoclam-zmp-backend/config"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common"
	"github.com/TruongHoang2004/ngoclam-zmp-backend/internal/common/log"
	"github.com/gin-gonic/gin"
)

func APIKeyChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		const prefix = "Apikey "
		webhookKey := config.AppConfig.WebhookApiKey

		if !strings.HasPrefix(authHeader, prefix) || strings.TrimPrefix(authHeader, prefix) != webhookKey {

			log.Error(c.Request.Context(), "Invalid API key")
			c.JSON(http.StatusUnauthorized, common.ConvertErrorToResponse(common.ErrUnauthorized(c.Request.Context())))
			c.Abort()
			return
		}

		c.Next()
	}
}
