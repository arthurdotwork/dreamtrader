package handler

import (
	"net/http"

	"github.com/arthurdotwork/dreamtrader/core/internal/request"
	"github.com/arthurdotwork/dreamtrader/core/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(registerService service.RegisterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		var req request.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := req.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := registerService.Register(ctx, req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}
