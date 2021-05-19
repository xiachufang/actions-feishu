package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/certifi/gocertifi"
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

	var caCerts *x509.CertPool
	if rootCAS, err := gocertifi.CACerts(); err == nil {
		caCerts = rootCAS
	} else {
		setOutput("Couldn't load CA Certificates")
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCerts,
			},
		},
	}
	resp, err := client.Post(webhook.Address, "application/json", bytes.NewBuffer(buf))
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				setOutput("close body error")
			}
		}()
	}
	if err != nil {
		setOutput(err.Error())
		return
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

// TemplateMessage 用于构造并处理处理模版请求消息的类型,不属于飞书定义的任何消息类型
type TemplateMessage struct {
	WebHook
	TemplatePath   string
	TemplateValues map[string]interface{}
}

// TemplateMessage implement Message
func (m *TemplateMessage) Send() {
	tmpl, err := template.ParseFiles(m.TemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	var rawData bytes.Buffer
	err = tmpl.Execute(&rawData, m.TemplateValues)
	if err != nil {
		log.Fatal(err)
	}

	var body map[string]interface{}
	err = json.Unmarshal(rawData.Bytes(), &body)
	if err != nil {
		log.Fatal(err)
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
	case "__template__":
		rawTemplateValues := os.Getenv("INPUT_MSG_TEMPLATE_VALUES")
		var templateValues map[string]interface{}
		err := json.Unmarshal([]byte(rawTemplateValues), &templateValues)
		if err != nil {
			log.Fatalf("hint: msg-template-values must be a valid JSON string\n\nReason: %s", err.Error())
		}
		templatePath := os.Getenv("INPUT_MSG_TEMPLATE_PATH")
		return &TemplateMessage{
			WebHook:        w,
			TemplatePath:   templatePath,
			TemplateValues: templateValues,
		}
	default:
		return NoopMessage
	}
}

func main() {
	msg := parseInput()
	msg.Send()
}
