package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	user_dtos "github.com/keyloom/web-api/dtos/user"
	"github.com/keyloom/web-api/entities"
)

type UserController struct{}

var _ core.Controller = (*UserController)(nil)

func (uc *UserController) RegisterRoutes(engine *gin.Engine) {
	userGroup := engine.Group("/users")
	{
		userGroup.POST("/", uc.CreateHandler)
	}
}

// @Summary Create a new user
// @Param body body user_dtos.CreateUserDTO true "User creation data"
// @Description Create a new user with the provided email and password
// @Accept json
// @Produce json
// @Success 200 {object} entities.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users/ [post]
func (uc *UserController) CreateHandler(c *gin.Context) {
	var dto user_dtos.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity := (&entities.User{}).CreateNew()
	entity.SetEmail(dto.Email)
	err := entity.SetPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = entity.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, entity)
}
