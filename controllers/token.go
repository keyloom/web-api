package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	token_dtos "github.com/keyloom/web-api/dtos/token"
	"github.com/keyloom/web-api/entities"
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

	// load user by email
	user := (&entities.User{}).LoadByEmail(req.Username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// verify password
	passMatch := user.CheckPassword(req.Password)
	if !passMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// generate token
	token, err := (&core.TokenService{}).GenerateToken(user.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// respond with token
	c.JSON(http.StatusOK, token)
}
