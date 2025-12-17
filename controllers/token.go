package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	token_dtos "github.com/keyloom/web-api/dtos/token"
)

type TokenController struct{}

var _ core.Controller = (*TokenController)(nil)

func (tc *TokenController) RegisterRoutes(engine *gin.Engine) {
	tokenGroup := engine.Group("/tokens")
	{
		tokenGroup.POST("/", tc.TokenDispatchHandler)
	}
}

func (tc *TokenController) TokenDispatchHandler(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	if grantType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grant_type is required"})
	}

	switch grantType {
	case core.ClientCredentialsGrant:
		{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
		}

	case core.PasswordGrant:
		{
			tc.PasswordGrantHandler(c)
		}

	default:
		{
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported grant_type"})
		}
	}
}

func (tc *TokenController) PasswordGrantHandler(c *gin.Context) {
	var req token_dtos.PasswordGrantRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Validate user credentials here
}
