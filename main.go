package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/TastefulElk/clai/command"
	"github.com/TastefulElk/clai/logger"
	"github.com/TastefulElk/clai/openai"

	"github.com/integrii/flaggy"
)

var defaultModel string = "gpt-4o"

func main() {
	flaggy.SetName("GPT CLI")
	flaggy.SetDescription("Ask the CLI for what you want do do then confirm the result to run it.")

	var query string
	flaggy.String(&query, "q", "query", "The query you want to ask the CLI.")

	var model string
	flaggy.String(&model, "m", "model", "The model you want to use. (Default gpt-4o)")

	var apiToken string
	flaggy.String(&apiToken, "t", "api-token", "The API token you want to use.")

	var verbose bool
	flaggy.Bool(&verbose, "v", "verbose", "Enable verbose logging.")

	flaggy.Parse()

	log, cleanup := logger.GetLogger(verbose)
	defer cleanup()

	log.Println("parsed 'query': ", query)
	log.Println("parsed 'model': ", model)
	log.Println("parsed 'verbose': ", verbose)

	if strings.TrimSpace(query) == "" {
		log.Println("no query provided, exiting.")
		flaggy.ShowHelpAndExit("'query' is required.")
	}

	if strings.TrimSpace(apiToken) == "" {
		log.Println("no api token provided, checking environment var CLAI_OPENAI_TOKEN.")
		apiToken = os.Getenv("CLAI_OPENAI_TOKEN")
		if apiToken == "" {
			log.Println("no api token found in environment either, exiting.")
			flaggy.ShowHelpAndExit("No API token found in environment. Either pass it as a flag or set the CLAI_OPENAI_TOKEN environment variable.")
		}
	}

	if strings.TrimSpace(model) == "" {
		log.Println("no model provided, using default model ", defaultModel)
		model = defaultModel
	}

	shell, _ := command.GetShell()
	log.Println("detected shell: ", shell)

	res, err := openai.Query(openai.QueryArg{
		Query:    query,
		Model:    model,
		ApiToken: apiToken,
		Shell:    shell,
	})
	if err != nil {
		fmt.Println("Error generating result:", err)
		log.Fatal("Error querying:", err)
	}

	log.Println("suggested command(s): ", res)

	const (
		Green = "\033[32m"
		Reset = "\033[0m"
	)

	fmt.Println("Here's the suggested command(s) to run:")
	fmt.Println(Green + res + Reset)
	fmt.Println("Do you want to run it directly? y/N")

	// read user input
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		log.Println("user chose not to run the command.")
		return
	}

	log.Println("user chose to run the command - executing.")
	command.RunCommand(res)
	log.Println("command executed, exiting")
}
