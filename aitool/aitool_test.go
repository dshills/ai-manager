package aitool_test

import (
	"encoding/json"
	"testing"

	"github.com/dshills/ai-manager/aitool"
)

func TestFuncJSON(t *testing.T) {
	tool := aitool.NewTool("GetWeather", "get the current weather in a given location",
		aitool.NewString("location", "The city and state, e.g. San Francisco. CA", true),
		aitool.NewString("unit", "Temperature units celsius or fahrenheit", false, "celsius", "fahrenheit"))

	_, err := json.MarshalIndent(tool, "", "\t")
	if err != nil {
		t.Error(err)
	}
}

func TestFuncJSON2(t *testing.T) {
	str1 := aitool.NewString("location", "The city and state, e.g. San Francisco. CA", true)
	str2 := aitool.NewString("unit", "Temperature units celsius or fahrenheit", false, "celsius", "fahrenheit")
	intProp := aitool.NewNumeric("myint", "Int test", true, aitool.NumericInt)
	aryProp := aitool.NewArray("myarray", "", true, intProp)
	strctProp := aitool.NewStruct("mystruct", "Struct test", true, str1, str2, intProp)

	tool := aitool.NewTool("GetWeather", "get the current weather in a given location", str1, str2, aryProp, strctProp)

	_, err := json.MarshalIndent(tool, "", "\t")
	if err != nil {
		t.Error(err)
	}
}
