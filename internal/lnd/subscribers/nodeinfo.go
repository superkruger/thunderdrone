package lnd

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"log"
	"time"
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

	for {
		select {
		case <-time.After(10 * time.Second):
			{
				getInfoResp, err := nis.client.GetInfo(nis.context, &lnrpc.GetInfoRequest{})

				if err != nil {
					return err
				}

				log.Println(getInfoResp)
			}
		case <-nis.context.Done(): // will execute if cancel func is called.
			return nil
		}

	}
}
