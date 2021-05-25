package sms

type Response interface {
	IsOk() bool
	GetBody() map[string]interface{}
}

type Client interface {
	Send(message Message) (Response, error)
}

type Message interface {
	GetNumbers() []string
	GetBody() string
}

type Msg struct {
	numbers []string
	body string
}

func (m Msg) GetNumbers() []string  {

	return m.numbers
}

func (m Msg) GetBody() string  {

	return m.body
}

func NewMsg(n []string, b string) Message {

	return Msg{
		numbers: n,
		body:    b,
	}
}