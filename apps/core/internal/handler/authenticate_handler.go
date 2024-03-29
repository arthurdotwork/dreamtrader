package handler

import (
	"net/http"

	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthenticateHandler(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req request.AuthenticateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := authService.Authenticate(ctx, req)
		if err != nil {
			// do not return any content to avoid leaking information
			c.JSON(http.StatusInternalServerError, gin.H{})
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
