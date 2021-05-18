package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Message interface of message
type Message interface {
	Send()
}

// WebHook webhook address
type WebHook struct {
	Address string
}

// set actions output
func setOutput(output string) {
	fmt.Printf(`::set-output name=response::%s\n`, output)
}

// post to feishu
func (webhook *WebHook) post(body interface{}) {
	buf, err := json.Marshal(body)
	if err != nil {
		setOutput(err.Error())
		return
	}

	resp, err := http.Post(webhook.Address, "application/json", bytes.NewBuffer(buf))
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				setOutput("close body error")
			}
		}()
	}
	if err != nil {
		setOutput(err.Error())
	}

	outout, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		setOutput(fmt.Sprintf("read response error: %s", err.Error()))
		return
	}
	setOutput(string(outout))
}

// TextMessage 文本类型消息
type TextMessage struct {
	WebHook
	Text string
}

// Send implement Message
func (m *TextMessage) Send() {
	body := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": m.Text,
		},
	}
	m.post(body)
}

// PostMessage Post 类型消息
type PostMessage struct {
	WebHook
	Title   string
	Content string
}

// Send implement Message
func (m *PostMessage) Send() {
	body := map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": m.Title,
					"content": [][]map[string]interface{}{
						{{
							"tag":  "text",
							"text": m.Content,
						},
						},
					},
				},
			},
		},
	}
	m.post(body)
}

type noopMessage struct{}

func (m *noopMessage) Send() {}

// NoopMessage no-op message
var NoopMessage = &noopMessage{}

func parseInput() Message {
	webhook := os.Getenv("INPUT_WEBHOOK")
	if webhook == "" {
		setOutput("webhook required")
		return NoopMessage
	}

	w := WebHook{Address: webhook}
	msgType := os.Getenv("INPUT_MESSAGE_TYPE")
	switch msgType {
	case "post":
		return &PostMessage{
			WebHook: w,
			Title:   os.Getenv("INPUT_TITLE"),
			Content: os.Getenv("INPUT_CONTENT"),
		}
	case "text":
		return &TextMessage{
			WebHook: w,
			Text:    os.Getenv("INPUT_CONTENT"),
		}
	default:
		return NoopMessage
	}
}

func main() {
	msg := parseInput()
	msg.Send()
}
