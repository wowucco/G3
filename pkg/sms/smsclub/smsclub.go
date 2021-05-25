package smsclub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/wowucco/G3/pkg/sms"
	"io/ioutil"
	"net/http"
)

const apiSendUrl = "https://im.smsclub.mobi/sms/send"

type Client struct {
	token string
	from  string
}

type Message struct {
	From    string   `json:"src_addr"`
	Phone   []string `json:"phone"`
	Message string   `json:"message"`
}

type Config struct {
	Token string
	From  string
}

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

func NewClient(cfg Config) Client {

	return Client{
		token: cfg.Token,
		from:  cfg.From,
	}
}

func (c Client) Send(message sms.Message) (sms.Response, error) {

	var status sms.Response

	body, err := json.Marshal(Message{
		From:    c.from,
		Phone:   message.GetNumbers(),
		Message: message.GetBody(),
	})

	if err != nil {
		return status, err
	}

	req, err := http.NewRequest("POST", apiSendUrl, bytes.NewBuffer(body))

	if err != nil {
		return status, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return status, err
	}

	var jsn map[string]interface{}

	b, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(b, &jsn); err != nil {
		return status, err
	}

	defer res.Body.Close()

	status = Response{status: true, response: jsn}

	return status, nil
}
