package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	Speech string `json:"speech"`
	Party  string `json:"party"`
}

type RequestN8NWebhook struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

type N8NWebhookResponse struct {
	Code    int
	Message string
	Hint    string
}

type N8NEnvVars struct {
	WebhookUrl     string
	webhookTestUrl string
}

type GeneralEnvVars struct {
	Env string
}

type Config struct {
	N8NEnvVars     N8NEnvVars
	GeneralEnvVars GeneralEnvVars
}

// get the right webhook url depending on the enviroment
func (c *Config) n8nUrlToUse() string {
	if c.GeneralEnvVars.Env == "test" {
		return c.N8NEnvVars.webhookTestUrl
	}

	return c.N8NEnvVars.WebhookUrl
}

func getValidEnv(key string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		panic(fmt.Sprintf("empty env var: %s", key))
	}

	return envValue
}

func getEnvsVariables() *Config {
	n8nVars := N8NEnvVars{
		WebhookUrl:     getValidEnv("N8N_WEBHOOK_URL"),
		webhookTestUrl: getValidEnv("N8N_WEBHOOK_TEST_URL"),
	}

	generalVars := GeneralEnvVars{
		Env: getValidEnv("ENV"),
	}

	return &Config{
		N8NEnvVars:     n8nVars,
		GeneralEnvVars: generalVars,
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load env variables")
	}

	config := getEnvsVariables()
	n8nRequest := RequestN8NWebhook{
		Type: "speech",
		Message: Message{
			Party:  "au-treck-is-my",
			Speech: "uma pequena vila oriental da idade media, com uma feira acontecendo, algumas pessoas que estao comprando na feira estao olhando assustados na direcao de voces",
		},
	}

	var n8nRequestBuff bytes.Buffer
	err = json.NewEncoder(&n8nRequestBuff).Encode(n8nRequest)
	if err != nil {
		fmt.Printf("error to marshal n8nRequest: %v", err)
	}

	res, err := http.Post(config.n8nUrlToUse(), "application/json", &n8nRequestBuff)
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
