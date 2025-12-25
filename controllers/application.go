package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	application_dtos "github.com/keyloom/web-api/dtos/application"
	"github.com/keyloom/web-api/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationController struct{}

var _ core.Controller = (*ApplicationController)(nil)

func (ac *ApplicationController) RegisterRoutes(engine *gin.Engine) {
	appGroup := engine.Group("/applications")
	{
		appGroup.POST("/", ac.CreateHandler)
		appGroup.GET("/", ac.GetAllHandler)
		appGroup.GET("/:id", ac.GetByIDHandler)
		appGroup.PUT("/:id", ac.UpdateHandler)
	}
}

// @Summary Create a new application
// @Param body body application_dtos.CreateApplicationDTO true "Application creation data"
// @Description Create a new application with the provided name, description, and audience ID
// @Accept json
// @Produce json
// @Success 201 {object} entities.Application
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /applications/ [post]
// @Tags Applications
func (ac *ApplicationController) CreateHandler(c *gin.Context) {
	var dto application_dtos.CreateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	entity := (&entities.Application{}).CreateNew()
	entity.Name = dto.Name
	entity.Description = dto.Description
	audienceID, err := primitive.ObjectIDFromHex(dto.AudienceID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid audience ID"})
		return
	}
	audienceEntity := (&entities.Audience{}).LoadByID(dto.AudienceID)
	if audienceEntity == nil {
		c.JSON(404, gin.H{"error": "Audience not found"})
		return
	}
	entity.AudienceID = audienceID
	entity.Audience = *audienceEntity
	entity.ClientID = primitive.NewObjectID().Hex()

	entity.Save()
	c.JSON(201, entity)
}

// @Summary Get all applications with pagination
// @Param limit query int false "Number of applications to return" default(10)
// @Param page query int false "Page number" default(1)
// @Description Retrieve a paginated list of all applications
// @Accept json
// @Produce json
// @Success 200 {array} []entities.Application
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /applications/ [get]
// @Tags Applications
func (ac *ApplicationController) GetAllHandler(c *gin.Context) {
	topParam := c.DefaultQuery("limit", "10")
	pageParam := c.DefaultQuery("page", "1")

	top, err := strconv.Atoi(topParam)
	if err != nil || top <= 0 {
		top = 10
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}

	applicationEntity := &entities.Application{}
	applications := applicationEntity.LoadAll(top, page)
	if applications == nil {
		c.JSON(404, gin.H{"error": "No applications found"})
		return
	}
	c.JSON(200, applications)
}

// @Summary Get application by ID
// @Param id path string true "Application ID"
// @Description Retrieve an application by its ID
// @Accept json
// @Produce json
// @Success 200 {object} entities.Application
// @Failure 404 {object} interface{}
// @Router /applications/{id} [get]
// @Tags Applications
func (ac *ApplicationController) GetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	applicationEntity := (&entities.Application{}).LoadByID(id)
	if applicationEntity == nil {
		c.JSON(404, gin.H{"error": "Application not found"})
		return
	}
	c.JSON(200, applicationEntity)
}

// @Summary Update an existing application
// @Param id path string true "Application ID"
// @Param body body application_dtos.CreateApplicationDTO true "Application update data"
// @Description Update an existing application's name and description
// @Accept json
// @Produce json
// @Success 200 {object} entities.Application
// @Failure 400 {object} interface{}
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /applications/{id} [put]
// @Tags Applications
func (ac *ApplicationController) UpdateHandler(c *gin.Context) {
	id := c.Param("id")
	var dto application_dtos.CreateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	applicationEntity := (&entities.Application{}).LoadByID(id)
	if applicationEntity == nil {
		c.JSON(404, gin.H{"error": "Application not found"})
		return
	}
	applicationEntity.Name = dto.Name
	applicationEntity.Description = dto.Description
	err := applicationEntity.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update application"})
		return
	}
	c.JSON(200, applicationEntity)
}
