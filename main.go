package main

import (
	"example.com/go-gpt/openai"
	"fmt"
	"github.com/integrii/flaggy"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

	var shell string
	flaggy.String(&shell, "s", "shell", "The shall that you want the commands to be executable in.")

	var verbose bool
	flaggy.Bool(&verbose, "v", "verbose", "Enable verbose logging.")

	flaggy.Parse()

	if strings.TrimSpace(query) == "" {
		flaggy.ShowHelpAndExit("'query' is required.")
	}

	if apiToken == "" {
		apiToken = os.Getenv("CLAI_OPENAI_TOKEN")
		if apiToken == "" {
			flaggy.ShowHelpAndExit("No API token found in environment. Either pass it as a flag or set the CLAI_OPENAI_TOKEN environment variable.")
		}
	}

	if shell == "" {
		shell = os.Getenv("SHELL")
	}

	if model == "" {
		model = defaultModel
	}

	res, err := openai.Query(openai.QueryArg{
		Query:    query,
		Model:    model,
		ApiToken: apiToken,
		Shell:    shell,
	})

	if err != nil {
		log.Fatal("Error querying:", err)
	}

	const (
		Red    = "\033[31m"
		Green  = "\033[32m"
		Yellow = "\033[33m"
		Blue   = "\033[34m"
		Reset  = "\033[0m"
	)

	fmt.Println("Here's the suggested command(s) to run:")
	fmt.Println(Green + res + Reset)
	fmt.Println("Do you want to run it directly? y/N")

	// read user input
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		return
	}

	runCommand(res)
}

// getShell dynamically detects the shell being used
func getShell() (string, string) {
	if runtime.GOOS == "windows" {
		// Windows: prefer PowerShell if available, else fallback to cmd.exe
		if os.Getenv("ComSpec") != "" {
			return "cmd.exe", "/C"
		}
		return "powershell.exe", "-Command"
	} else {
		// Unix-like systems: use the SHELL environment variable or fallback to /bin/sh
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh" // fallback to /bin/sh if SHELL is not set
		}

		return shell, "-c"
	}
}

// runCommand runs the generated command in the detected shell
func runCommand(command string) error {
	shell, shellFlag := getShell()

	// Create the command to execute using the detected shell
	cmd := exec.Command(shell, shellFlag, command)

	// Set up to pipe the output and error directly to stdout/stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}
