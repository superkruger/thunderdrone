package services

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/superkruger/thunderdrone/internal/interfaces"
	"github.com/superkruger/thunderdrone/internal/repositories"
	"io"
	"log"
)

type ConnectionDetails struct {
	GRPCAddress       string
	TLSFileBytes      []byte
	MacaroonFileBytes []byte
}

type nodeSettingsService struct {
	repo      repositories.NodeSettingsRepo
	lndClient interfaces.LndClient
}

func NewNodeSettingsService(repo repositories.NodeSettingsRepo) NodeSettingsService {
	return &nodeSettingsService{
		repo: repo,
	}
}

func (nss *nodeSettingsService) SetConnectionDetails(localNode repositories.LocalNode) (repositories.LocalNode, error) {

	if localNode.TLSFile == nil {
		return localNode, fmt.Errorf("TLS file wasn't supplied")
	}

	if localNode.MacaroonFile == nil {
		return localNode, fmt.Errorf("macaroon file wasn't supplied")
	}

	localNode.TLSFileName = &localNode.TLSFile.Filename
	tlsDataFile, err := localNode.TLSFile.Open()
	if err != nil {
		return localNode, err
	}
	tlsData, err := io.ReadAll(tlsDataFile)
	if err != nil {
		return localNode, err
	}
	localNode.TLSDataBytes = tlsData

	localNode.MacaroonFileName = &localNode.MacaroonFile.Filename
	macaroonDataFile, err := localNode.MacaroonFile.Open()
	if err != nil {
		return localNode, err
	}
	macaroonData, err := io.ReadAll(macaroonDataFile)
	if err != nil {
		return localNode, err
	}
	localNode.MacaroonDataBytes = macaroonData

	if localNode.NodeId == "undefined" {
		_uuid, err := uuid.NewV4()
		if err != nil {
			return localNode, err
		}
		uuidStr := _uuid.String()
		localNode.NodeId = uuidStr

		log.Printf("creating local node %v\n", localNode)
		err = nss.repo.CreateLocalNode(localNode)
		if err != nil {
			return localNode, err
		}
	} else {
		log.Printf("updating local node %v\n", localNode)
		err = nss.repo.UpdateLocalNode(localNode)
	}

	if nss.lndClient != nil {
		err = nss.lndClient.Restart()
	}

	return localNode, err
}

func (nss *nodeSettingsService) GetConnectionDetails() (ConnectionDetails, error) {
	//localNodeDetails, err := nss.repo.GetLocalNode()
	//if err != nil {
	//	return ConnectionDetails{}, err
	//}
	//if (localNodeDetails.GRPCAddress == nil) || (localNodeDetails.TLSDataBytes == nil) || (localNodeDetails.
	//	MacaroonDataBytes == nil) {
	//	return ConnectionDetails{}, errors.New("missing node details")
	//}
	//fmt.Printf("GetConnectionDetails: %v\n", localNodeDetails)
	//return ConnectionDetails{
	//	GRPCAddress:       *localNodeDetails.GRPCAddress,
	//	TLSFileBytes:      localNodeDetails.TLSDataBytes,
	//	MacaroonFileBytes: localNodeDetails.MacaroonDataBytes}, nil

	return ConnectionDetails{}, nil
}

func (nss *nodeSettingsService) AllNodes() ([]repositories.LocalNode, error) {
	return nss.repo.GetLocalNodes()
}

func (nss *nodeSettingsService) SetPubKey(localNode repositories.LocalNode) error {
	return nss.repo.UpdatePubKey(localNode)
}

func (nss *nodeSettingsService) SetLndClient(client interfaces.LndClient) {
	nss.lndClient = client
}
