package lnd

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lightningnetwork/lnd/lnrpc"
	"time"
)

var localNodePubKeys []string

func InitLocalNodesList(ctx context.Context, client lnrpc.LightningClient, db *sqlx.DB) error {

	var pubKey *string
	var grpcAddress string

	q := `select grpc_address, pub_key from local_node;`
	r, err := db.Query(q)

	for r.Next() {
		err := r.Scan(&grpcAddress, &pubKey)
		if err != nil {
			return errors.Wrapf(err, "r.Scan(&grpcAddress, &pubKey)")
		}

		// If the pub key is missing from the local_node table, add it.
		if pubKey == nil || len(*pubKey) == 0 {
			pubKey, err = addMissingLocalPubkey(ctx, client, grpcAddress, db)
			if err != nil {
				return errors.Wrapf(err, "addMissingLocalPubkey(ctx, client, grpcAddress, db)")
			}
		}
		localNodePubKeys = append(localNodePubKeys, *pubKey)
	}
	if err != nil {
		return errors.Wrapf(err, "db.Query(%s)", q)
	}

	return nil
}

func addMissingLocalPubkey(ctx context.Context, client lnrpc.LightningClient, grpcAddress string,
	db *sqlx.DB) (r *string, err error) {

	// Get the public key of our node
	ni, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return nil, errors.Wrapf(err, "client.GetNodeInfo(ctx, &lnrpc.GetInfoRequest{})")
	}

	const q = `update local_node set(pub_key, updated_on) = ($1, $2) where grpc_address = $3`

	_, err = db.Exec(q,
		ni.IdentityPubkey,
		time.Now().UTC(),
		grpcAddress,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "tx.Exec(%v, %v, %v, %v)",
			q,
			ni.IdentityPubkey,
			time.Now().UTC(),
			grpcAddress,
		)
	}

	return &ni.IdentityPubkey, nil
}
