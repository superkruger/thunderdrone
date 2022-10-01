package lnd

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	lnd "github.com/superkruger/thunderdrone/internal/lnd/consumer"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Subscribe(ctx context.Context, conn *grpc.ClientConn, db *sqlx.DB, localNodeId int) error {

	client := lnrpc.NewLightningClient(conn)

	errs, ctx := errgroup.WithContext(ctx)

	// Store a list of public keys belonging to our nodes
	err := lnd.InitLocalNodesList(ctx, client, db)
	if err != nil {
		return err
	}

	err = errs.Wait()
	fmt.Println("Subscriptions all ended")

	return err
}
