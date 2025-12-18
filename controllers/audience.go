package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	audience_dtos "github.com/keyloom/web-api/dtos/audience"
	"github.com/keyloom/web-api/entities"
)

type AudienceController struct{}

var _ core.Controller = (*AudienceController)(nil)

func (ac *AudienceController) RegisterRoutes(engine *gin.Engine) {
	audienceGroup := engine.Group("/audiences")
	{
		audienceGroup.POST("/", ac.CreateHandler)
	}
}

// @Summary Create a new audience
// @Param body body audience_dtos.CreateAudienceDTO true "Audience creation data"
// @Description Create a new audience with the provided display name and description
// @Accept json
// @Produce json
// @Success 201 {object} entities.Audience
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /audiences/ [post]
func (ac *AudienceController) CreateHandler(c *gin.Context) {
	var dto audience_dtos.CreateAudienceDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	entity := (&entities.Audience{}).CreateNew()
	entity.DisplayName = dto.DisplayName
	entity.Description = dto.Description

	// Generate a unique name (slug) from the display name
	// Replace spaces with hyphens and convert to lowercase
	entity.Name = strings.ToLower(strings.ReplaceAll(dto.DisplayName, " ", "-"))
	err := entity.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create audience"})
		return
	}
	c.JSON(201, entity)
}
