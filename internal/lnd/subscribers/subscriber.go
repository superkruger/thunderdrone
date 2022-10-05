package lnd

type Subscriber interface {
	Subscribe() error
}
