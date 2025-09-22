package main

import (
	"testing"
)

type TestCase struct {
	input    string
	expected []string
}

func TestCleanInput(t *testing.T) {
	cases := []TestCase{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "wuS poppin JIMbo",
			expected: []string{"wus", "poppin", "jimbo"},
		},
		// add more cases here
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		expected := c.expected
		if len(actual) != len(expected) {
			t.Errorf("input slice length does not match expected slice length, FAIL")
		}
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("input word does not equal expected word, FAIL")
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
