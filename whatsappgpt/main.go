package main

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int      `json:"index"`
	Message struct { // daria para usar a struct de message de cima, mas é um exemplo de utilização
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func GenerateGPTText(query string) (string, error) {
	req := Request{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens: 150,
	}

	reqJson, err := json.Marshal(req)

	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.newReader(reqJson))

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer YOUR_API_KEY")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close() // Roda tudo que ta em baixo do defer e depois roda o defer

}
