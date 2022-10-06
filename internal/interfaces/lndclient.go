package interfaces

type LndClient interface {
	Start() error
	Restart() error
}
