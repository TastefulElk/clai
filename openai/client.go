package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var url = "https://api.openai.com/v1/chat/completions"

type QueryArg struct {
	Query    string
	Model    string
	ApiToken string
	Shell    string
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

func Query(queryArgs QueryArg) (string, error) {
	systemMessage := `
    You are a helpful assistant living in a CLI application.
    Your job is to respond to the users problem with a CLI command or a series of commands that solves the problem posed.

    Send the raw commands, no formatting.

    If you don't understand the problem, or for any reason are unable to respond with cli commands, send the string "n/a".
    NEVER return with anything else other than commands that are directly executable or "n/a".

    The shell for which the commands should be executed in will be passed as the first part of the query

    Ex. input: 'bash: I want to find all files in the current directory that contains the text 'foo' in them.'
    Ex. output: 'grep -r 'foo' .'

    Ex. input: 'powershell: I want to tail the log file log.txt'
    Ex. output: 'tail -f log.txt'

    Your responses will, after confirmation, be executed directly, so be careful with what you output.
  `

	payload := Payload{
		Model: queryArgs.Model,
		Messages: []message{
			{
				Role:    "system",
				Content: systemMessage,
			},
			{
				Role:    "user",
				Content: queryArgs.Shell + ": " + queryArgs.Query,
			},
		},
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error marshalling request body:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+queryArgs.ApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// Parse response
	var chatResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(responseBody, &chatResponse)
	if err != nil {
		log.Fatal("Error parsing response:", err)
	}

	commands := chatResponse.Choices[0].Message.Content

	if commands == "n/a" {
		return "", errors.New("no result")
	}

	return commands, nil
}
