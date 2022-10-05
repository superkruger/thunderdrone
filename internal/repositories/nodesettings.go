package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"mime/multipart"
	"time"
)

type LocalNode struct {
	NodeId            *string               `json:"localNodeId" form:"localNodeId" db:"local_node_id"`
	GRPCAddress       *string               `json:"grpcAddress" form:"grpcAddress" db:"grpc_address"`
	TLSFileName       *string               `json:"tlsFileName" db:"tls_file_name"`
	TLSFile           *multipart.FileHeader `form:"tlsFile"`
	TLSDataBytes      []byte                `db:"tls_data"`
	MacaroonFileName  *string               `json:"macaroonFileName" db:"macaroon_file_name"`
	MacaroonFile      *multipart.FileHeader `form:"macaroonFile"`
	MacaroonDataBytes []byte                `db:"macaroon_data"`
	CreateAt          time.Time             `json:"createdAt" db:"created_at"`
	UpdatedAt         *time.Time            `json:"updatedAt"  db:"updated_at"`
	PubKey            *string               `db:"pub_key"`
}

type NodeSettingsRepo interface {
	GetLocalNode(nodeId string) (LocalNode, error)
	UpdateLocalNode(node LocalNode) error
	CreateLocalNode(node LocalNode) error
	DeleteLocalNode(nodeId string) error
	GetLocalNodes() ([]LocalNode, error)
	UpdatePubKey(node LocalNode) error
}

type nodeSettingsRepo struct {
	db *sqlx.DB
}

func NewNodeSettingsRepo(db *sqlx.DB) NodeSettingsRepo {
	return &nodeSettingsRepo{
		db: db,
	}
}

func (nsr *nodeSettingsRepo) GetLocalNode(nodeId string) (localNodeData LocalNode, err error) {
	err = nsr.db.Get(&localNodeData, `
SELECT
  grpc_address,
  tls_data,
  macaroon_data,
  pub_key
FROM local_node
WHERE 
    local_node_id = $1
LIMIT 1;`, nodeId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LocalNode{}, nil
		}
		return LocalNode{}, errors.Wrap(err, "Unable to execute SQL query")
	}
	return localNodeData, nil
}

func (nsr *nodeSettingsRepo) UpdateLocalNode(node LocalNode) (err error) {
	_, err = nsr.db.Exec(`
UPDATE local_node SET
  grpc_address = $1,
  tls_file_name = $2,
  tls_data = $3,
  macaroon_file_name = $4,
  macaroon_data = $5,
  updated_on = $6
WHERE 
  local_node_id = $7;
`, node.GRPCAddress, node.TLSFileName, node.TLSDataBytes, node.MacaroonFileName, node.MacaroonDataBytes, time.Now().UTC(), node.NodeId)
	if err != nil {
		return errors.Wrap(err, "Unable to Update Local Node Details")
	}
	return nil
}

func (nsr *nodeSettingsRepo) CreateLocalNode(node LocalNode) (err error) {
	_, err = nsr.db.Exec(`
INSERT into local_node ('local_node_id', 'grpc_address', 'tls_file_name', 'tls_data', 'macaroon_file_name', 'macaroon_data', 'created_at')
VALUES ($1, $2, $3, $4, $5, $6, $7);
`, node.NodeId, node.GRPCAddress, node.TLSFileName, node.TLSDataBytes, node.MacaroonFileName, node.MacaroonDataBytes, time.Now().UTC())
	if err != nil {
		return errors.Wrap(err, "Unable to Create Local Node")
	}
	return nil
}

func (nsr *nodeSettingsRepo) DeleteLocalNode(nodeId string) (err error) {
	_, err = nsr.db.Exec(`
DELETE from local_node
WHERE local_node_id = $1;
`, nodeId)
	if err != nil {
		return errors.Wrap(err, "Unable to Delete Local Node")
	}
	return nil
}

func (nsr *nodeSettingsRepo) GetLocalNodes() (localNodes []LocalNode, err error) {
	err = nsr.db.Select(&localNodes, `
SELECT
  grpc_address,
  tls_data,
  macaroon_data,
  pub_key
FROM local_node;`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []LocalNode{}, nil
		}
		return []LocalNode{}, errors.Wrap(err, "Unable to execute SQL query")
	}
	return localNodes, nil
}

func (nsr *nodeSettingsRepo) UpdatePubKey(node LocalNode) (err error) {
	_, err = nsr.db.Exec(`
UPDATE local_node SET
  pub_key = $1
  updated_on = $2
WHERE 
	local_node_id = $3;
`, node.PubKey, time.Now().UTC(), node.NodeId)
	if err != nil {
		return errors.Wrap(err, "Unable to Update Local Node Details")
	}
	return nil
}
