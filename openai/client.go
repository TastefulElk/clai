package openai

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

var url = "https://api.openai.com/v1/chat/completions"

type QueryArg struct {
	Query    string
	Model    string
	ApiToken string
}

func Query(queryArgs QueryArg) string {
	systemMessage := `
    You are a helpful assistant living in a CLI application.
    Your job is to respond to the users problem with a CLI command or a series of commands that solves the problem posed.

    Ex. input: "I want to find all files in the current directory that contains the text 'foo' in them."
    Ex. output: "grep -r 'foo' ."

    Ex. input: "I want to tail the log file log.txt"
    Ex. output: "tail -f log.txt"

    Your responses will, after confirmation, be executed directly, so be careful with what you output.
  `

	body := `{
    "model": "` + queryArgs.Model + `",
    "messages": [
      {
        "role": "system",
        "content": "` + systemMessage + `"
      },
      {
        "role": "user",
        "content": "` + queryArgs.Query + `"
      }
    ]
  }`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))

	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+queryArgs.ApiToken)

	fmt.Println("Sending request: ", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)

	responseBody := new(bytes.Buffer)

	_, err = responseBody.ReadFrom(resp.Body)

	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	log.Println("Response Body:", responseBody)
	return responseBody.String()
}
