package contact

type IRecallForm interface {
	GetPhone() string
	GetMessage() string
}

type IBuyOnClickForm interface {
	GetPhone() string
	GetProductId() int
}
