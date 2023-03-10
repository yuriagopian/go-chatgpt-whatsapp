package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	// "io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(reqJson))

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

	// respBody, err := ioutil.ReadAll(response.Body)
	respBody, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func parseBase64RequestData(r string) (string, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(r)

	if err != nil {
		return "", err
	}
	data, err := url.ParseQuery(string(dataBytes))

	if data.Has("Body") {
		return data.Get("Body"), nil
	}

	return "", errors.New("Body not found")
}
func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, err := parseBase64RequestData(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	text, err := GenerateGPTText(result)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       text,
	}, nil
}

func main() {
	lambda.Start(process)
}
