package controllers

import (
	"strconv"
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
		audienceGroup.GET("/", ac.GetAllHandler)
		audienceGroup.GET("/:id", ac.GetByIDHandler)
		audienceGroup.PUT("/:id", ac.UpdateHandler)
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
// @Tags Audiences
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

// @Summary Get all audiences with pagination
// @Param limit query int false "Number of audiences to return" default(10)
// @Param page query int false "Page number" default(1)
// @Description Retrieve a paginated list of audiences
// @Accept json
// @Produce json
// @Success 200 {array} entities.Audience
// @Failure 500 {object} interface{}
// @Router /audiences/ [get]
// @Tags Audiences
func (ac *AudienceController) GetAllHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")

	// Convert limit and page to integers
	top, err := strconv.Atoi(limit)
	if err != nil || top <= 0 {
		top = 10
	}
	pg, err := strconv.Atoi(page)
	if err != nil || pg <= 0 {
		pg = 1
	}
	audience := &entities.Audience{}
	audiences := audience.LoadAll(top, pg)
	c.JSON(200, audiences)
}

// @Summary Get audience by ID
// @Param id path string true "Audience ID"
// @Description Retrieve an audience by its ID
// @Accept json
// @Produce json
// @Success 200 {object} entities.Audience
// @Failure 404 {object} interface{}
// @Router /audiences/{id} [get]
// @Tags Audiences
func (ac *AudienceController) GetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	audience := (&entities.Audience{}).LoadByID(id)
	if audience == nil {
		c.JSON(404, gin.H{"error": "Audience not found"})
		return
	}
	c.JSON(200, audience)
}

// @Summary Update an existing audience
// @Param id path string true "Audience ID"
// @Param body body audience_dtos.UpdateAudienceDTO true "Audience update data"
// @Description Update an existing audience's display name and description
// @Accept json
// @Produce json
// @Success 200 {object} entities.Audience
// @Failure 400 {object} interface{}
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /audiences/{id} [put]
// @Tags Audiences
func (ac *AudienceController) UpdateHandler(c *gin.Context) {
	dto := audience_dtos.UpdateAudienceDTO{}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	id := c.Param("id")
	audience := (&entities.Audience{}).LoadByID(id)
	if audience == nil {
		c.JSON(404, gin.H{"error": "Audience not found"})
		return
	}
	audience.DisplayName = dto.DisplayName
	audience.Description = dto.Description
	err := audience.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update audience"})
		return
	}
	c.JSON(200, audience)
}
