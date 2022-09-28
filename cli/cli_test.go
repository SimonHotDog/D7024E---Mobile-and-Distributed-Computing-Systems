package cli

import (
	"fmt"
	"testing"
)

func TestParseCommand(t *testing.T) {
	var tests = []struct {
		input    string
		expected command
	}{
		{"help", command{"help", ""}},
		{"get myhash", command{"get", "myhash"}},
		{"put my message", command{"put", "my message"}},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("Parse \"%s\"", test.input)
		t.Run(testname, func(t *testing.T) {
			actual := parseCommand(test.input)
			if actual != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, actual)
			}
		})
	}
}

func TestRemoveDoubleQuotes(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"\"\"", ""},
		{"\"hello world", "\"hello world"},
		{"hello world\"", "hello world\""},
		{"\"hello world\"", "hello world"},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("Remove double quotes from '%s'", test.input)
		t.Run(testname, func(t *testing.T) {
			actual := removeDoubleQuotes(test.input)
			if actual != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, actual)
			}
		})
	}
}
