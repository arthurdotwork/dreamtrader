package handler

import (
	"net/http"

	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/gin-gonic/gin"
)

func RefreshAuthenticationHandler(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if len(c.GetHeader("Authorization")) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		accessToken, refreshToken, err := authService.RefreshAuthentication(ctx, c.GetHeader("Authorization")[7:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": gin.H{
				"token":      accessToken.AccessToken,
				"expires_at": accessToken.ExpiresAt,
			},
			"refresh_token": gin.H{
				"token":      refreshToken.RefreshToken,
				"expires_at": refreshToken.ExpiresAt,
			},
		})
	}
}
