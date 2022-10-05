package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"github.com/superkruger/thunderdrone/internal/services"
	"net/http"
)

type nodeSettingsRoutes struct {
	nodeSettings services.NodeSettingsService
}

func NewNodeSettingsRoutes(nodeSettings services.NodeSettingsService) Routable {
	return &nodeSettingsRoutes{
		nodeSettings: nodeSettings,
	}
}

func (nsr *nodeSettingsRoutes) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("nodesettings", func(c *gin.Context) { nsr.UpdateNodeSettings(c) })
}

func (nsr *nodeSettingsRoutes) UpdateNodeSettings(c *gin.Context) {

	var localNode repositories.LocalNode

	if err := c.Bind(&localNode); err != nil {
		c.AbortWithError(400, err)
		return
	}

	localNode, err := nsr.nodeSettings.SetConnectionDetails(localNode)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.JSON(http.StatusOK, localNode)
}
