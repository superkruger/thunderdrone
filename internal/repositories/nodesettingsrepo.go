package repositories

type NodeSettingsRepo interface {
	GetLocalNode(nodeId string) (LocalNode, error)
	UpdateLocalNode(node LocalNode) error
	CreateLocalNode(node LocalNode) error
	DeleteLocalNode(nodeId string) error
	GetLocalNodes() ([]LocalNode, error)
	UpdatePubKey(node LocalNode) error
}
