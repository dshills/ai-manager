# ai-manager

Interact with all the modern AIs

## Example Usage

```go
package main

import (
	"ai-manager/ai"
	"ai-manager/aigen/gemini"
	"ai-manager/aigen/mistral"
	"ai-manager/aigen/openai"
	"ai-manager/aiserial"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	serialPath = ".ai-manager"
	openAIKey  = "<YOUR OpenAI API Key>"
	geminiKey  = "<YOUR Gemini API Key>"
	mistralKey = "<YOUR Mistral API Key>"
)

var openAIModels = []string{"gpt-4", "gpt-3.5-turbo", "gpt-4-turbo-preview"}
var geminiModels = []string{"gemini-pro", "gemini-1.0-pro-latest", "gemini-1.0-pro"}
var mistralModels = []string{"mistral-small-latest", "mistral-medium-latest", "mistral-large-latest"}

func main() {
	home := os.Getenv("HOME")
	path := filepath.Join(home, serialPath)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		log.Printf("[FATAL] %v", err)
		os.Exit(1)
	}
	serializer := aiserial.New(path)
	aimgr := ai.New(serializer)

	aimgr.RegisterGenerator(openai.AIName, openAIKey, openAIModels, openai.Generator)
	aimgr.RegisterGenerator(gemini.AIName, geminiKey, geminiModels, gemini.Generator)
	aimgr.RegisterGenerator(mistral.AIName, mistralKey, mistralModels, mistral.Generator)

	aimgr.NewThread(openai.AIName, "gpt-3.5-turbo")

	output := make(chan string)

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
