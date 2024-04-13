package aitool_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dshills/ai-manager/aitool"
)

// Test functions
func testFunc1(x int, y string) bool {
	return x > 0 && y != ""
}

func testFunc2(a, b float64, c uint64, d uint) float64 {
	return a + b + float64(c) + float64(d)
}

func getCurrentWeather(location, unit string) string {
	return fmt.Sprintf("%v %v\n", location, unit)
}

func TestGetFunctionSchema(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		function    interface{}
		expected    string
		paramInfo   []string
	}{
		{
			name:        "testFunc1",
			description: "Generic test function 1",
			function:    testFunc1,
		},
		{
			name:        "testFunc2",
			description: "Generic test function 1",
			function:    testFunc2,
		},
		{
			name:        "getCurrentWeather",
			description: "Get the current weather in a given location",
			function:    getCurrentWeather,
			paramInfo:   []string{"location:The city and state, e.g. San Francisco, CA", "units:celsius or fahrenheit"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool, err := aitool.ToolFromFunc(tc.name, tc.description, tc.function, tc.paramInfo...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			js, err := json.MarshalIndent(tool, "", "  ")
			if err != nil {
				t.Error(err)
				return
			}

			fmt.Println(string(js))
		})
	}
}
