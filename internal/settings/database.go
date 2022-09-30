package settings

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func getLocalNodeConnectionDetails(db *sqlx.DB, implementation string) (localNodeData LocalNode, err error) {
	err = db.Get(&localNodeData, `
SELECT
  grpc_address,
  tls_data,
  macaroon_data
FROM local_node
WHERE 
    implementation = $1
LIMIT 1;`, implementation)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LocalNode{}, nil
		}
		return LocalNode{}, errors.Wrap(err, "Unable to execute SQL query")
	}
	return localNodeData, nil
}

func updateLocalNodeDetails(db *sqlx.DB, node LocalNode) (err error) {
	_, err = db.Exec(`
UPDATE local_node SET
  grpc_address = $1,
  tls_file_name = $2,
  tls_data = $3,
  macaroon_file_name = $4,
  macaroon_data = $5,
  updated_on = $6
WHERE 
  implementation = $7;
`, node.GRPCAddress, node.TLSFileName, node.TLSDataBytes, node.MacaroonFileName, node.MacaroonDataBytes, time.Now().UTC(), node.Implementation)
	if err != nil {
		return errors.Wrap(err, "Unable to Update Local Node Details")
	}
	return nil
}
