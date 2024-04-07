# ai-manager

Interact with all the modern AIs

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
		{AIName: openaiName, Model: gpt4, APIKey: openaiKey, BaseURL: openaiBaseURL, Generator: openai.Generator},
		{AIName: openaiName, Model: gpt35turbo, APIKey: openaiKey, BaseURL: openaiBaseURL, Generator: openai.Generator},
		{AIName: geminiName, Model: gemini1pro, APIKey: geminiKey, BaseURL: geminiBaseURL, Generator: gemini.Generator},
	}

	// Register the models
	aimgr.RegisterGenerators(models...)

	thread := ai.ThreadData{
		AIName: openaiName,
		Model:  gpt4,
	}
	// Create a thread to converse with
	if err := aimgr.NewThread(thread); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a channel to receive data
	output := make(chan string)

	// Start chatting
	go aimgr.CurrentThread().Generate(output, "Write a story about a superhero cat named Bitty")

	for {
		msg := <-output
		switch msg {
		case ai.ErrorStart:
			fmt.Println("Received an error!")
			continue
		case ai.ResponseStart:
			continue

		case ai.ErrorComplete:
			os.Exit(1)
		case ai.ResponseComplete:
			os.Exit(0)
		}
		fmt.Println(msg)
	}
}
```
