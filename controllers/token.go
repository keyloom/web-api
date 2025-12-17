package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
)

type TokenController struct{}

var _ core.Controller = (*TokenController)(nil)

func (tc *TokenController) RegisterRoutes(engine *gin.Engine) {
	tokenGroup := engine.Group("/tokens")
	{
		tokenGroup.POST("/", tc.PostDispatchHandler)
	}
}

func (tc *TokenController) PostDispatchHandler(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	if grantType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "grant_type is required"})
	}

	switch grantType {
	case core.ClientCredentialsGrant:
		{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
		}
	}
}
