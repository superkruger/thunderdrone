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
}

func NewNodeInfoScubscriber(client lnrpc.LightningClient, localNode repositories.LocalNode) Subscriber {
	return &nodeInfoSubscriber{
		client:    client,
		localNode: localNode,
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
