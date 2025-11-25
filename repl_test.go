package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Mama give me PICKachoo",
			expected: []string{"mama", "give", "me", "pickachoo"},
		},
	// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) = %v; want %v", c.input, actual, c.expected)
			t.Fail()
		}
		// Check the length of the actual slice against the expected slice
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) = %v; want %v", c.input, actual, c.expected)
				t.Fail()
			}
			// Check each word in the slice
		}
	}

}