package lnd

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"gopkg.in/macaroon.v2"
	"io/ioutil"
	"path"
	"time"
)

func Start() error {

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(info, warn, err))

	tlsCertPath := path.Join("/app/tls.cert")
	macaroonPath := path.Join("/app/admin.macaroon")

	tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
	if err != nil {
		fmt.Println("Cannot get node tls credentials", err)
		return err
	}

	//	spew.Println("TLS", tlsCreds)

	macaroonBytes, err := ioutil.ReadFile(macaroonPath)
	if err != nil {
		fmt.Println("Cannot read macaroon file", err)
		return err
	}

	//	spew.Println("Macaroon", macaroonBytes)

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		fmt.Println("Cannot unmarshal macaroon", err)
		return err
	}

	macCred, err := macaroons.NewMacaroonCredential(mac)
	if err != nil {
		fmt.Println("Cannot get macaroon credentials", err)
		return err
	}
	opts := []grpc.DialOption{
		grpc.WithReturnConnectionError(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(tlsCreds),
		grpc.WithPerRPCCredentials(macCred),
	}

	fmt.Println("dialing")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "Alice:10009", opts...) //grpc.Dial("localhost:10009", opts...)
	if err != nil {
		fmt.Println("cannot dial to lnd", err)
		return err
	}
	fmt.Println("getting client")
	client := lnrpc.NewLightningClient(conn)

	fmt.Println("get info")

	ctx = context.Background()
	getInfoResp, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		fmt.Println("Cannot get info from node:", err)
		return err
	}
	spew.Dump(getInfoResp)

	return nil
}
