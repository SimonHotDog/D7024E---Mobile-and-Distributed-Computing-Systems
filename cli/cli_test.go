package cli

import (
	"bytes"
	"d7024e/kademlia"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNewCli(t *testing.T) {
	testname := "Test if cli is created"
	t.Run(testname, func(t *testing.T) {
		kademlia := &kademlia.Kademlia{}
		cli := NewCli(os.Stdin, os.Stdout, kademlia)

		if cli == nil {
			t.Errorf("Expected cli, got nil")
			t.FailNow()
		}

		if cli.context != kademlia {
			t.Errorf("Expected kademlia, got nil")
			t.Fail()
		}

		if cli.in == nil {
			t.Errorf("Expected stdin, got nil")
			t.Fail()
		}

		if cli.out == nil {
			t.Errorf("Expected stdout, got nil")
			t.Fail()
		}
	})
}

func TestOpenGreeting(t *testing.T) {
	testname := "Test if greeting is printed"
	t.Run(testname, func(t *testing.T) {
		kademlia := &kademlia.Kademlia{}
		var outBuffer bytes.Buffer
		readBuffer := bytes.NewBufferString("EOF")
		cli := NewCli(&outBuffer, readBuffer, kademlia)

		cli.Open(true)
		actual := outBuffer.String()

		if strings.Contains(actual, "Hello") == false {
			t.Errorf("Expected greeting, got %s", actual)
		}
	})
}

func TestNoInput(t *testing.T) {
	testname := "Test if no input is ignored"
	t.Run(testname, func(t *testing.T) {
		input := "\nEOF"
		expected := ""
		kademlia := &kademlia.Kademlia{}
		var outBuffer bytes.Buffer
		readBuffer := bytes.NewBufferString(input)
		cli := NewCli(&outBuffer, readBuffer, kademlia)

		cli.Open(false)
		actual := outBuffer.String()

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	})
}

func TestUnknownCommand(t *testing.T) {
	testname := "Test if unknown command is interpreted as unknown"
	t.Run(testname, func(t *testing.T) {
		input := "Uknonwn command\nEOF"
		expectedToContain := "not found"
		kademlia := &kademlia.Kademlia{}
		var outBuffer bytes.Buffer
		readBuffer := bytes.NewBufferString(input)
		cli := NewCli(&outBuffer, readBuffer, kademlia)

		cli.Open(false)
		actual := outBuffer.String()

		if strings.Contains(actual, expectedToContain) == false {
			t.Errorf("Expected to contain %s, got %s", expectedToContain, actual)
		}
	})
}

func TestACommand(t *testing.T) {
	testname := "Test if a command is executed"
	t.Run(testname, func(t *testing.T) {
		input := "help\nEOF"
		kademlia := &kademlia.Kademlia{}
		var outBuffer bytes.Buffer
		readBuffer := bytes.NewBufferString(input)
		cli := NewCli(&outBuffer, readBuffer, kademlia)

		cli.Open(false)
		actual := outBuffer.String()

		if len(actual) == 0 {
			t.Errorf("Expected to any output, got noting")
		}
	})
}

func TestSplitCommand(t *testing.T) {
	var tests = []struct {
		input        string
		expectedName string
		expectedArgs string
	}{
		{"", "", ""},
		{"help", "help", ""},
		{"get myhash", "get", "myhash"},
		{"put my message", "put", "my message"},
		{"put my\n message\n", "put", "my\n message"},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("Parse \"%s\"", test.input)
		t.Run(testname, func(t *testing.T) {
			actualName, actualArgs := splitCommand(test.input)
			if actualName != test.expectedName || actualArgs != test.expectedArgs {
				t.Errorf("Expected (%s, %s) got (%s, %s)", actualName, actualArgs, test.expectedName, test.expectedArgs)
			}
		})
	}
}
