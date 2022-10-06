package lnd

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"log"
	"time"
)

type peerInfoSubscriber struct {
	client    lnrpc.LightningClient
	localNode repositories.LocalNode
	context   context.Context
}

func NewPeerInfoSubscriber(client lnrpc.LightningClient, localNode repositories.LocalNode, ctx context.Context) Subscriber {
	return &peerInfoSubscriber{
		client:    client,
		localNode: localNode,
		context:   ctx,
	}
}

func (s *peerInfoSubscriber) Subscribe() error {

	for {
		select {
		case <-time.After(10 * time.Second):
			{
				peers, err := s.client.ListPeers(s.context, &lnrpc.ListPeersRequest{})

				if err != nil {
					return err
				}

				log.Println("Peers", peers)
			}
		case <-s.context.Done(): // will execute if cancel func is called.
			return nil
		}

	}
}
