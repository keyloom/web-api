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
	tokenGroup := engine.Group("/token")
	{
		tokenGroup.POST("/", tc.TokenDispatchHandler)
		tokenGroup.GET("/me/validate", tc.ValidateToken)
	}
}

// @Summary Token dispatch endpoint
// @Param grant_type formData string true "Grant type"
// @Param username formData string false "Username for password grant"
// @Param password formData string false "Password for password grant"
// @Param client_id formData string false "Client ID for client credentials grant"
// @Description Dispatch tokens based on the provided grant type
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} token_dtos.AccessTokenResponse
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /token/ [post]
// @Tags Tokens
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
	token, err := (&core.TokenService{}).GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// respond with token
	c.JSON(http.StatusOK, token)
}

// @Summary Validate token endpoint
// @Description Validate the provided JWT token
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} interface{}
// @Failure 401 {object} interface{}
// @Router /token/me/validate [get]
// @Tags Tokens
// @Security ApiKeyAuth
func (tc *TokenController) ValidateToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	tokenString = tokenString[len("Bearer "):]

	token, err := (&core.TokenService{}).ValidateToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token is valid",
		"payload": token,
	})
}
