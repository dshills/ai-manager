package openai

import (
	"fmt"
	"os"
	"testing"

	"github.com/dshills/ai-manager/ai"
	"github.com/dshills/ai-manager/aitool"
)

func TestFuncCall(t *testing.T) {
	const (
		model   = "gpt-3.5-turbo"
		baseURL = "https://api.openai.com/v1"
	)
	apiKey := os.Getenv("OPENAI_APIKEY")
	if apiKey == "" {
		t.Fatal("OPENAI_APIKEY not found")
	}

	msg := ai.Message{Role: "user", Text: "What is the weather in Des Moines today?"}
	conv := ai.Conversation{msg}
	locProp := aitool.NewString("location", "The city and state e.g. San Francisco, CA", true)
	unitProp := aitool.NewString("unit", "units to return", true, "celsius", "fahrenheit")
	tool := aitool.NewTool("GetCurrentWeather", "Get the current weather in a given location", locProp, unitProp)
	tools := []aitool.Tool{*tool}

	gen := Generator{model: model, apiKey: apiKey, baseURL: baseURL, tools: make(map[string]aitool.Tool)}

	/*
		req := gen.NewRequest(conv, nil, tools)
		js, err := json.MarshalIndent(&req, "", "\t")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(js))
	*/

	resp, err := gen.Generate(conv, nil, tools)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
}
