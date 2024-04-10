# ai-manager

Interact with all the modern AIs

## Supported Models

Anthropic: claude-3-haiku-20240307
Anthropic: claude-3-opus-20240229
Anthropic: claude-3-sonnet-20240229
Gemini: gemini-1.0-pro
Gemini: gemini-1.0-pro-latest
Gemini: gemini-pro
Mistral: mistral-large-latest
Mistral: mistral-medium-latest
Mistral: mistral-small-latest
Ollama: Any
OpenAI: gpt-3.5-turbo
OpenAI: gpt-4
OpenAI: gpt-4-turbo-preview

Others can be added by writing to the Generator interface.

## Example Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/dshills/ai-manager/ai"
	"github.com/dshills/ai-manager/aigen/gemini"
	"github.com/dshills/ai-manager/aigen/openai"
)

const (
	openaiName    = "OpenAI"
	openaiKey     = "<YOUR OpenAI API Key>"
	gpt4          = "gpt-4"
	gpt35turbo    = "gpt-3.5-turbo"
	openaiBaseURL = "https://api.openai.com/v1"
)
const (
	geminiName    = "Gemini"
	geminiKey     = "<YOUR Gemini API Key>"
	gemini1pro    = "gemini-1.0-pro"
	geminiBaseURL = "https://generativelanguage.googleapis.com/v1beta"
)

func main() {
	// Create the manager
	aimgr := ai.New()

	// Models we want to use
	models := []ai.Model{
		{AIName: openaiName, Model: gpt4, APIKey: openaiKey, BaseURL: openaiBaseURL, Generator: openai.New()},
		{AIName: openaiName, Model: gpt35turbo, APIKey: openaiKey, BaseURL: openaiBaseURL, Generator: openai.New()},
		{AIName: geminiName, Model: gemini1pro, APIKey: geminiKey, BaseURL: geminiBaseURL, Generator: gemini.New()},
	}

	// Register the models
	err := aimgr.RegisterGenerators(models...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	threadData := ai.ThreadData{
		AIName: openaiName,
		Model:  gpt4,
	}
	// Create a thread to converse with
	thread, err := aimgr.NewThread(threadData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	qry := "Write a story about a superhero cat named Bitty"

	resp, err := thread.Converse(qry)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("AI said: %s\n", resp.Message.Text)
	fmt.Printf("Response time: %v\n", resp.Elapsed)
	fmt.Printf("Token Cose: %v\n", resp.Usage.TotalTokens)
}
```
