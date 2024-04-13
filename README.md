# ai-manager

Interact with all the modern AIs

**Important Note:** The API/interface is changing rapidly as I add new functionality and capabilities. I am also refactoring out redundant or confusing structures. If you are using the package and this is causing issues let me know and I will take a more careful approach to refactoring.

## Supported Models

- Anthropic: claude-3-haiku-20240307
- Anthropic: claude-3-opus-20240229
- Anthropic: claude-3-sonnet-20240229
- Gemini: gemini-1.0-pro
- Gemini: gemini-1.0-pro-latest
- Gemini: gemini-pro
- Mistral: mistral-large-latest
- Mistral: mistral-medium-latest
- Mistral: mistral-small-latest
- Ollama: Any
- OpenAI: gpt-3.5-turbo
- OpenAI: gpt-4
- OpenAI: gpt-4-turbo-preview

Others can be added by writing to the Generator interface.

## Basic Example

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
	openaiKey       = "<YOUR OpenAI API Key>"
	modelGPT4       = "gpt-4"
	modelGPT35Turbo = "gpt-3.5-turbo"
	openaiBaseURL   = "https://api.openai.com/v1"
)
const (
	geminiKey       = "<YOUR Gemini API Key>"
	modelGemini1Pro = "gemini-1.0-pro"
	geminiBaseURL   = "https://generativelanguage.googleapis.com/v1beta"
)

func main() {
	// Create the manager
	aimgr := ai.New()

	genGPT4 := openai.New(modelGPT4, openaiKey, openaiBaseURL)
	genGPT35Turbo := openai.New(modelGPT35Turbo, openaiKey, openaiBaseURL)
	genGemini1Pro := gemini.New(modelGemini1Pro, geminiKey, geminiBaseURL)

	// Register the models
	err := aimgr.RegisterGenerators(genGPT35Turbo, genGPT4, genGemini1Pro)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a thread to converse with
	thread, err := aimgr.NewThread(ai.NewThreadData(modelGPT35Turbo))
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

## Function Calling Example
```go
package main

import (
	"fmt"
	"os"

	"github.com/dshills/ai-manager/ai"
	"github.com/dshills/ai-manager/aigen/openai"
	"github.com/dshills/ai-manager/aitool"
)

const (
	openaiKey       = "<YOUR OpenAI API KEY>"
	modelGPT35Turbo = "gpt-3.5-turbo"
	openaiBaseURL   = "https://api.openai.com/v1"
)

func GetCurrentWeather(location, unit string) string {
	/// Simulate calling an external system for weather data
	return fmt.Sprintf("%s %s", location, unit)
}

func main() {
	// Create the manager
	aimgr := ai.New()

	// Create a generator
	genGPT35Turbo := openai.New(modelGPT35Turbo, openaiKey, openaiBaseURL)

	// Register the generator
	err := aimgr.RegisterGenerators(genGPT35Turbo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	threadData := ai.NewThreadData(modelGPT35Turbo)
	weatherTool, err := aitool.ToolFromFunc("GetCurrentWeather", "Get the current weather iin a given location",
		GetCurrentWeather, "location:The city and state. e.g. San Francisco, CA", "unit:celsius or fahrenheit")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Set our tool function
	threadData.Tools = []aitool.Tool{*weatherTool}
	// Set the temperature to 1 it's the default anyway
	threadData.MetaData = []ai.Meta{{Key: "temperature", Value: "1"}}

	// Create a thread to converse with
	thread, err := aimgr.NewThread(threadData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	qry := "What is the weather in West Des Moines, Iowa"

	resp, err := thread.Converse(qry)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, toolCall := range resp.ToolCalls {
		fcall, err := toolCall.FuncString()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(fcall)
	}
}
```
Output: GetCurrentWeather("West Des Moines, Iowa", "celsius")

