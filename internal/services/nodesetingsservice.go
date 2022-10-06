package services

import (
	"github.com/superkruger/thunderdrone/internal/interfaces"
	"github.com/superkruger/thunderdrone/internal/repositories"
)

type NodeSettingsService interface {
	GetConnectionDetails() (ConnectionDetails, error)
	SetConnectionDetails(node repositories.LocalNode) (repositories.LocalNode, error)
	AllNodes() ([]repositories.LocalNode, error)
	SetPubKey(localNode repositories.LocalNode) error
	SetLndClient(client interfaces.LndClient)
}
