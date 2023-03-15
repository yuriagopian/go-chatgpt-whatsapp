package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	// "io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
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

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
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

	chatGptApikey := os.Getenv("CHAT_GPT_API_KEY")

	fmt.Println(chatGptApikey)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer"+chatGptApikey)
	fmt.Println(request.Header)
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
	data, _ := url.ParseQuery(string(dataBytes))

	if data.Has("Body") {
		return data.Get("Body"), nil
	}

	return "", errors.New("body not found")
}
func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("request", request)
	result, err := parseBase64RequestData(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	println("aqui2 ")
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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	fmt.Println("1 aqui")
	fmt.Println(godotenv.Load(".env"))
	fmt.Println(os.Getenv("TWILIO_SID"))
	lambda.Start(Handler)
	fmt.Println("teste")
}
