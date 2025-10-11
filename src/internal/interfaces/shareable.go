package interfaces

type Shareable interface {
	GetID() int64
	GetType() string
	GetTitle() string
}
