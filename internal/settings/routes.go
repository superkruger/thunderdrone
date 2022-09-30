package settings

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func RegisterSettingRoutes(r *gin.RouterGroup, db *sqlx.DB) {
	r.POST("nodesettings", func(c *gin.Context) { updateNodeSettings(c, db) })
}

type LocalNode struct {
	NodeId            int                   `json:"localNodeId" form:"localNodeId" db:"local_node_id"`
	Implementation    string                `json:"implementation" form:"implementation" db:"implementation"`
	GRPCAddress       *string               `json:"grpcAddress" form:"grpcAddress" db:"grpc_address"`
	TLSFileName       *string               `json:"tlsFileName" db:"tls_file_name"`
	TLSFile           *multipart.FileHeader `form:"tlsFile"`
	MacaroonFileName  *string               `json:"macaroonFileName" db:"macaroon_file_name"`
	MacaroonFile      *multipart.FileHeader `form:"macaroonFile"`
	CreateOn          time.Time             `json:"createdOn" db:"created_on"`
	UpdatedOn         *time.Time            `json:"updatedOn"  db:"updated_on"`
	TLSDataBytes      []byte                `db:"tls_data"`
	MacaroonDataBytes []byte                `db:"macaroon_data"`
}

func updateNodeSettings(c *gin.Context, db *sqlx.DB) {

	var localNode LocalNode

	if err := c.Bind(&localNode); err != nil {
		c.AbortWithError(400, err)
		return
	}
	localNode.NodeId = 1

	if localNode.TLSFile == nil {
		c.AbortWithError(400, errors.New("TLS file wasn't supplied"))
		return
	}

	if localNode.MacaroonFile == nil {
		c.AbortWithError(400, errors.New("macaroon file wasn't supplied"))
		return
	}

	localNode.TLSFileName = &localNode.TLSFile.Filename
	tlsDataFile, err := localNode.TLSFile.Open()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	tlsData, err := io.ReadAll(tlsDataFile)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	localNode.TLSDataBytes = tlsData

	localNode.MacaroonFileName = &localNode.MacaroonFile.Filename
	macaroonDataFile, err := localNode.MacaroonFile.Open()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	macaroonData, err := io.ReadAll(macaroonDataFile)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	localNode.MacaroonDataBytes = macaroonData

	err = updateLocalNodeDetails(db, localNode)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusOK, localNode)
}
