package handler

import (
	"net/http"

	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/gin-gonic/gin"
)

func VerifyAuthenticationHandler(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if len(c.GetHeader("Authorization")) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if err := authService.Verify(ctx, c.GetHeader("Authorization")[7:]); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
