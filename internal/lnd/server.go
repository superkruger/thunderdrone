package lnd

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/superkruger/thunderdrone/internal/settings"
	"google.golang.org/grpc/grpclog"
	"log"
)

func Start(ctx context.Context, db *sqlx.DB) error {

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(info, warn, err))

	connectionDetails, err := settings.GetConnectionDetails(db)
	if err != nil && err.Error() != "Missing node details" {
		fmt.Printf("failed to get node connection details: %v", err)
		return err
	}

	conn, err := Connect(
		connectionDetails.GRPCAddress,
		connectionDetails.TLSFileBytes,
		connectionDetails.MacaroonFileBytes)
	if err != nil {
		log.Printf("Failed to connect to lnd: %v\n", err)
		return err
	}

	err = Subscribe(ctx, conn, db, 1)
	if err != nil {
		log.Printf("Failed to subscribe to lnd: %v\n", err)
		return err
	}

	return nil
}
