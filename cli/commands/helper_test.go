package commands

import (
	"fmt"
	"testing"
)

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
			actual := RemoveDoubleQuotes(test.input)
			if actual != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, actual)
			}
		})
	}
}

func TestGetAllCommands(t *testing.T) {
	testname := "Test if commands are returned"
	t.Run(testname, func(t *testing.T) {
		actual := AllCommands()
		if len(actual) == 0 {
			t.Errorf("Expected commands, got none")
		}
	})
}
