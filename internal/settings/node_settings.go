package settings

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type connectionDetails struct {
	GRPCAddress       string
	TLSFileBytes      []byte
	MacaroonFileBytes []byte
}

func GetConnectionDetails(db *sqlx.DB) (connectionDetails, error) {
	localNodeDetails, err := getLocalNodeConnectionDetails(db, "LND")
	if err != nil {
		return connectionDetails{}, err
	}
	if (localNodeDetails.GRPCAddress == nil) || (localNodeDetails.TLSDataBytes == nil) || (localNodeDetails.
		MacaroonDataBytes == nil) {
		return connectionDetails{}, errors.New("missing node details")
	}
	fmt.Printf("GetConnectionDetails: %v\n", localNodeDetails)
	return connectionDetails{
		GRPCAddress:       *localNodeDetails.GRPCAddress,
		TLSFileBytes:      localNodeDetails.TLSDataBytes,
		MacaroonFileBytes: localNodeDetails.MacaroonDataBytes}, nil
}
