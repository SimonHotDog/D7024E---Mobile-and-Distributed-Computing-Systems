package util

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestGetNextCommand(t *testing.T) {
	var tests = []struct {
		increments uint64
		expected   int64
	}{
		{1, 1},
		{1000, 1000},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("Increment counter %d times", test.increments)
		t.Run(testname, func(t *testing.T) {

			c := MakeCounter()
			for i := uint64(1); i < test.increments; i++ {
				go c.GetNext()
			}
			time.Sleep(1 * time.Second)
			actual := c.GetNext()

			if actual != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, actual)
			}
		})
	}
}

func TestGetNextWraparoundCommand(t *testing.T) {
	expected := int64(500)
	testname := "Test wraparound when max of int64 is reached"
	t.Run(testname, func(t *testing.T) {

		c := MakeCounter()
		c.i = math.MaxInt64 - 500

		for i := 1; i < 1000; i++ {
			go c.GetNext()
		}
		time.Sleep(1 * time.Second)
		actual := c.GetNext()

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})
}

func TestIncreaseCommand(t *testing.T) {
	var tests = []struct {
		increments uint64
		expected   int64
	}{
		{1, 1},
		{1000, 1000},
	}
	for _, test := range tests {
		testname := fmt.Sprintf("Increment counter %d times", test.increments)
		t.Run(testname, func(t *testing.T) {

			c := MakeCounter()
			for i := uint64(0); i < test.increments; i++ {
				go c.Increase()
			}
			time.Sleep(1 * time.Second)
			actual := c.GetNext() - 1

			if actual != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, actual)
			}
		})
	}
}

func TestIncreaseWraparoundCommand(t *testing.T) {
	expected := int64(500)
	testname := "Test wraparound when max of int64 is reached"
	t.Run(testname, func(t *testing.T) {

		c := MakeCounter()
		c.i = math.MaxInt64 - 500

		for i := 0; i < 1000; i++ {
			go c.Increase()
		}
		time.Sleep(1 * time.Second)
		actual := c.GetNext() - 1

		if actual != expected {
			t.Errorf("Expected %d, got %d", expected, actual)
		}
	})
}
