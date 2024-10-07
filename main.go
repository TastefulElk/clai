package main

import (
	"example.com/go-gpt/openai"
	"fmt"
	"log"

	"github.com/integrii/flaggy"
)

func main() {
	flaggy.SetName("GPT CLI")
	flaggy.SetDescription("Ask the CLI for what you want do do then confirm the result to run it.")

	var query string
	flaggy.String(&query, "q", "query", "The query you want to ask the CLI.")

	var model string
	flaggy.String(&model, "m", "model", "The model you want to use.")

	var apiToken string
	flaggy.String(&apiToken, "t", "api-token", "The API token you want to use.")

	flaggy.Parse()

	log.Print("Query: ", query)
	log.Print("Model: ", model)
	log.Print("API Token: ", apiToken)

	fmt.Println("Querying...")
	res := openai.Query(openai.QueryArg{
		Query:    query,
		Model:    model,
		ApiToken: apiToken,
	})

	fmt.Println("Result: ", res)
}
