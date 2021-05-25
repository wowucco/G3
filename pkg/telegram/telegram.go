package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ParseModeMarkdown = "Markdown"

type Config struct {
	ApiUrl string
	BotId  string
}

func NewMsg(chat, body, parseMode string) Message {

	return Msg{
		chat:      chat,
		body:      body,
		parseMode: parseMode,
	}
}

type Msg struct {
	chat      string
	body      string
	parseMode string
}

func (m Msg) GetChat() string {

	return m.chat
}

func (m Msg) GetBody() string {

	return m.body
}

func (m Msg) GetParseMode() string {

	if m.parseMode == "" {
		return ParseModeMarkdown
	} else {
		return m.parseMode
	}
}

type Resp struct {
	status bool
	body map[string]interface{}
}

func (r Resp) IsOk() bool {

	return r.status
}

func (r Resp) GetBody() map[string]interface{} {

	return r.body
}

type Body struct {
	Text string `json:"text"`
	Chat string `json:"chat_id"`
	Mode string `json:"parse_mode"`
}

func NewTelegramClient(cfg Config) (Client, error) {

	if cfg.ApiUrl == "" || cfg.BotId == "" {
		return nil, errors.New("telegram: failed create client, miss api url or bot id")
	}

	return &Telegram{
		apiUrl:     cfg.ApiUrl,
		botId:      cfg.BotId,
		httpClient: &http.Client{},
	}, nil
}

type Telegram struct {
	apiUrl     string
	botId      string
	httpClient *http.Client
}

func (t Telegram) Send(msg Message) (Response, error) {

	var status Response

	body, err := json.Marshal(Body{
		Text: msg.GetBody(),
		Chat: msg.GetChat(),
		Mode: msg.GetParseMode(),
	})

	if err != nil {
		return status, err
	}

	req, err := http.NewRequest("POST", t.url(), bytes.NewBuffer(body))

	if err != nil {
		return status, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := t.httpClient.Do(req)

	if err != nil {
		return status, err
	}

	defer res.Body.Close()

	var jsn map[string]interface{}

	b, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(b, &jsn); err != nil {
		return status, err
	}

	status = Resp{
		status: true,
		body:   jsn,
	}

	return status, nil
}

func (t Telegram) url() string {

	url := t.apiUrl

	lst := url[len(url) - 1:]

	if lst[0] != '/' {
		url = fmt.Sprintf("%s/", url)
	}

	return fmt.Sprintf("%sbot%s/sendMessage", url, t.botId)
}

type Mock struct {

}

func NewMock() (Client, error) {

	return Mock{}, nil
}

func (m Mock) Send(msg Message) (Response, error) {

	return Resp{
		status: true,
		body:   map[string]interface{}{
			"status": "ok",
			"message": fmt.Sprintf("[mock telegram message] text: %v, chat: %v", msg.GetBody(), msg.GetChat()),
		},
	}, nil
}

// Contracts
type Message interface {
	GetChat() string
	GetBody() string
	GetParseMode() string
}

type Response interface {
	IsOk() bool
	GetBody() map[string]interface{}
}

type Client interface {
	Send(msg Message) (Response, error)
}
