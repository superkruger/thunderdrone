package lnd

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/superkruger/thunderdrone/internal/repositories"
)

type nodeInfoSubscriber struct {
	client    lnrpc.LightningClient
	localNode repositories.LocalNode
	context   context.Context
}

func NewNodeInfoScubscriber(client lnrpc.LightningClient, localNode repositories.LocalNode, ctx context.Context) Subscriber {
	return &nodeInfoSubscriber{
		client:    client,
		localNode: localNode,
		context:   ctx,
	}
}

func (nis *nodeInfoSubscriber) Subscribe() error {

	ctx := context.Background()
	getInfoResp, err := nis.client.GetInfo(ctx, &lnrpc.GetInfoRequest{})

	if err != nil {
		return err
	}

	fmt.Println(getInfoResp)

	return nil
}
