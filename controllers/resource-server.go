package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keyloom/web-api/core"
	"github.com/keyloom/web-api/dtos/resource_server_dtos"
	"github.com/keyloom/web-api/entities"
)

type ResourceServerController struct{}

var _ core.Controller = (*ResourceServerController)(nil)

func (ac *ResourceServerController) RegisterRoutes(engine *gin.Engine) {
	resourceServerGroup := engine.Group("/resource-servers")
	{
		resourceServerGroup.POST("/", ac.CreateHandler)
		resourceServerGroup.GET("/", ac.GetAllHandler)
		resourceServerGroup.GET("/:id", ac.GetByIDHandler)
		resourceServerGroup.PUT("/:id", ac.UpdateHandler)
	}
}

// @Summary Create a new resource server
// @Param body body resource_server_dtos.CreateResourceServerDTO true "Resource server creation data"
// @Description Create a new resource server with the provided display name and description
// @Accept json
// @Produce json
// @Success 201 {object} entities.ResourceServer
// @Failure 400 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /resource-servers/ [post]
// @Tags ResourceServers
func (ac *ResourceServerController) CreateHandler(c *gin.Context) {
	var dto resource_server_dtos.CreateResourceServerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	entity := (&entities.ResourceServer{}).CreateNew()
	entity.DisplayName = dto.DisplayName
	entity.Description = dto.Description

	// Generate a unique name (slug) from the display name
	// Replace spaces with hyphens and convert to lowercase
	entity.Name = strings.ToLower(strings.ReplaceAll(dto.DisplayName, " ", "-"))
	err := entity.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create resource server"})
		return
	}
	c.JSON(201, entity)
}

// @Summary Get all resource servers with pagination
// @Param limit query int false "Number of resource servers to return" default(10)
// @Param page query int false "Page number" default(1)
// @Description Retrieve a paginated list of resource servers
// @Accept json
// @Produce json
// @Success 200 {array} entities.ResourceServer
// @Failure 500 {object} interface{}
// @Router /resource-servers/ [get]
// @Tags ResourceServers
func (ac *ResourceServerController) GetAllHandler(c *gin.Context) {
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
	resourceServer := &entities.ResourceServer{}
	resourceServers := resourceServer.LoadAll(top, pg)
	c.JSON(200, resourceServers)
}

// @Summary Get resource server by ID
// @Param id path string true "Resource Server ID"
// @Description Retrieve a resource server by its ID
// @Accept json
// @Produce json
// @Success 200 {object} entities.ResourceServer
// @Failure 404 {object} interface{}
// @Router /resource-servers/{id} [get]
// @Tags ResourceServers
func (ac *ResourceServerController) GetByIDHandler(c *gin.Context) {
	id := c.Param("id")
	resourceServer := (&entities.ResourceServer{}).LoadByID(id)
	if resourceServer == nil {
		c.JSON(404, gin.H{"error": "Resource Server not found"})
		return
	}
	c.JSON(200, resourceServer)
}

// @Summary Update an existing resource server
// @Param id path string true "Resource Server ID"
// @Param body body resource_server_dtos.UpdateResourceServerDTO true "Resource Server update data"
// @Description Update an existing resource server's display name and description
// @Accept json
// @Produce json
// @Success 200 {object} entities.ResourceServer
// @Failure 400 {object} interface{}
// @Failure 404 {object} interface{}
// @Failure 500 {object} interface{}
// @Router /resource-servers/{id} [put]
// @Tags ResourceServers
func (ac *ResourceServerController) UpdateHandler(c *gin.Context) {
	dto := resource_server_dtos.UpdateResourceServerDTO{}
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	id := c.Param("id")
	resourceServer := (&entities.ResourceServer{}).LoadByID(id)
	if resourceServer == nil {
		c.JSON(404, gin.H{"error": "Resource Server not found"})
		return
	}
	resourceServer.DisplayName = dto.DisplayName
	resourceServer.Description = dto.Description
	err := resourceServer.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update resource server"})
		return
	}
	c.JSON(200, resourceServer)
}
