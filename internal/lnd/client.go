package lnd

import (
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/lightningnetwork/lnd/lnrpc"
	lnd "github.com/superkruger/thunderdrone/internal/lnd/subscribers"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"github.com/superkruger/thunderdrone/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	"sync"
)

type LndClient interface {
	Start() error
}

type lndClient struct {
	context          context.Context
	localNodePubKeys []string
	nodeSettings     services.NodeSettingsService
}

func NewLndClient(context context.Context, nodeSettings services.NodeSettingsService) LndClient {
	return &lndClient{
		context:      context,
		nodeSettings: nodeSettings,
	}
}

func (lc *lndClient) Start() error {

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(info, warn, err))

	localNodes, err := lc.nodeSettings.AllNodes()
	if err != nil {
		return err
	}

	var wg = sync.WaitGroup{}

	for _, localNode := range localNodes {

		if (localNode.GRPCAddress == nil) || (localNode.TLSDataBytes == nil) || (localNode.MacaroonDataBytes == nil) {
			return fmt.Errorf("missing node connection details for %v", localNode.NodeId)
		}

		conn, err := Connect(
			*localNode.GRPCAddress,
			localNode.TLSDataBytes,
			localNode.MacaroonDataBytes)
		if err != nil {
			log.Printf("Failed to connect to lnd: %v\n", err)
			return err
		}

		wg.Add(1)
		go func(context context.Context, localNode repositories.LocalNode, conn grpc.ClientConnInterface) {

			err = lc.subscribe(localNode, conn)
			if err != nil {
				log.Printf("Failed to subscribe to lnd: %v\n", err)
			}
			log.Println("Node subscription done")
			wg.Done()
		}(lc.context, localNode, conn)
	}

	wg.Wait()

	return nil
}

func (lc *lndClient) subscribe(localNode repositories.LocalNode, conn grpc.ClientConnInterface) error {

	client := lnrpc.NewLightningClient(conn)

	// Initialise
	err := lc.initLocalNode(client, localNode)
	if err != nil {
		return err
	}

	var wg = sync.WaitGroup{}

	for _, subscriber := range lc.allSubscribers(localNode, client, lc.context) {
		wg.Add(1)
		go func(subscriber lnd.Subscriber) {

			err = subscriber.Subscribe()
			if err != nil {
				log.Println("Could not start subscriber")
			}
			log.Println("Subscriber done")
			wg.Done()
		}(subscriber)
	}

	wg.Wait()

	log.Println("Subscriptions all ended")

	return err
}

func (lc *lndClient) allSubscribers(localNode repositories.LocalNode, client lnrpc.LightningClient, ctx context.Context) []lnd.Subscriber {
	return []lnd.Subscriber{
		lnd.NewNodeInfoScubscriber(client, localNode, ctx),
	}
}

func (lc *lndClient) initLocalNode(client lnrpc.LightningClient, localNode repositories.LocalNode) error {

	if localNode.PubKey == nil || len(*localNode.PubKey) == 0 {
		pubKey, err := lc.addMissingLocalPubkey(localNode, client)
		if err != nil {
			return errors.Wrapf(err, "addMissingLocalPubkey")
		}
		lc.localNodePubKeys = append(lc.localNodePubKeys, *pubKey)
	}

	return nil
}

func (lc *lndClient) addMissingLocalPubkey(localNode repositories.LocalNode, client lnrpc.LightningClient) (r *string, err error) {

	// Get the public key of our node
	ni, err := client.GetInfo(lc.context, &lnrpc.GetInfoRequest{})
	if err != nil {
		return nil, errors.Wrapf(err, "client.GetNodeInfo(ctx, &lnrpc.GetInfoRequest{})")
	}

	localNode.PubKey = &ni.IdentityPubkey

	err = lc.nodeSettings.SetPubKey(localNode)
	if err != nil {
		return nil, err
	}

	return &ni.IdentityPubkey, nil
}
