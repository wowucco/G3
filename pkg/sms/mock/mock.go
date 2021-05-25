package mock

import (
	"fmt"
	"github.com/wowucco/G3/pkg/sms"
)

type Response struct {
	status bool
	response map[string]interface{}
}

func (r Response) IsOk() bool {
	return r.status == true
}

func (r Response) GetBody() map[string]interface{} {
	return r.response
}

type Client struct {}

func NewClient() Client {

	return Client{}
}

func (c Client) Send(message sms.Message) (sms.Response, error) {

	res := map[string]interface{}{
		"status": "ok",
		"message": fmt.Sprintf("mock sms to numbers: %v, text: %v", message.GetNumbers(), message.GetBody()),
	}

	return Response{status: true, response: res}, nil
}
