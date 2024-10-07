# CLAI

CLAI is a command line utility that'll generate, and optionally, execute
CLI commands for you based on your input.

## Installation

````bash
```sh
go install github.com/tastefulelk/clai
````

## Requirements

You need an OpenAI API key, which can be obtained for free by following the instructions [here](https://platform.openai.com/docs/quickstart/step-2-setup-your-api-key).

## Usage

First, you need to set your OpenAI API key as an environment variable:

```bash
export CLAI_OPENAI_KEY=your-key
```

Then run clai:

```bash
$ clai -q "Find all rows in log.txt that contains the words 'error' or 'warn'"
Here's the suggested command(s) to run:
grep -Ei 'error|warn' log.txt
Do you want to run it directly? y/N
```

### Options

- `-m|--model`: Specify the model to use. (Default `gpt-4o`)
- `-h|--help`: Show help message.
- `-v|--verbose`: Show verbose output.
- `-t|--api-token`: Specify OpenAI API token. This will take priority over the `CLAI_OPENAI_KEY` environment variable.
- `--version`: Show clai version

## Disclaimer

The user is solely responsible for reviewing and executing the generated commands. Please carefully review each command before running it, especially commands that modify system settings, files, or data. The creators of this tool are not liable for any damage or unintended consequences resulting from the execution of these commands.

## Disclaimer 2

I've never written a single of go before, don't judge me 🙈
