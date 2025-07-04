package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Speech string
}

type RequestN8NWebhook struct {
	Type    string
	Message Message
}

type N8NWebhookResponse struct {
	Code    int
	Message string
	Hint    string
}

func main() {
	n8nRequest := RequestN8NWebhook{
		Type: "speech",
		Message: Message{
			Speech: "caio recebe uma maça do mago",
		},
	}

	var n8nRequestBuff bytes.Buffer
	err := json.NewEncoder(&n8nRequestBuff).Encode(n8nRequest)
	if err != nil {
		fmt.Printf("error to marshal n8nRequest: %v", err)
	}

	res, err := http.Post("http://localhost:5678/webhook-test/f9969ca9-4bec-499a-94a5-c8bf796e3af5", "application/json", &n8nRequestBuff)
	if err != nil {
		fmt.Printf("error to call webhook: %s", err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error to read res body: %v", err)
	}

	var parsedResBody N8NWebhookResponse
	if err = json.Unmarshal(resBody, &parsedResBody); err != nil {
		fmt.Printf("error to unmarshal response: %v", err)
	}

	fmt.Println(parsedResBody.Message)

	if res.StatusCode != 200 {
		panic("webhook has return error")
	}
}
