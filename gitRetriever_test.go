package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeduplicate(t *testing.T) {

	tests := []struct {
		name           string
		usernames      []string
		expectedOutput []string
	}{
		{
			name: "no-duplicates",
			usernames: []string{
				"conor-wilson",
				"Daffy-Duck",
				"buggsBunny123",
			},
			expectedOutput: []string{
				"conor-wilson",
				"Daffy-Duck",
				"buggsBunny123",
			},
		},
		{
			name: "duplicates",
			usernames: []string{
				"conor-wilson",
				"conor-wilson",
				"Daffy-Duck",
			},
			expectedOutput: []string{
				"conor-wilson",
				"Daffy-Duck",
			},
		},
		{
			name:           "nil-usernames",
			usernames:      nil,
			expectedOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call function under test
			output := deduplicate(tt.usernames)

			// Confirm the output is as expected
			if len(output) != len(tt.expectedOutput) {
				t.Errorf("deduplicate() returned a slice of unexpected length:\nWant: %v\nGot:  %v", tt.expectedOutput, output)
			}
			for i := 0; i < len(tt.expectedOutput); i++ {
				if output[i] != tt.expectedOutput[i] {
					t.Errorf("deduplicate() returned an unexpected output:\nWant: %v\nGot:  %v", tt.expectedOutput, output)
				}
			}
		})
	}
}

func TestAlphabetiseUsers(t *testing.T) {

	tests := []struct {
		name           string
		users          []*User
		expectedOutput []*User
	}{
		{
			name: "happy-path",
			users: []*User{
				{
					Name: "Conor Wilson",
				},
				{
					Name: "Daffy Duck",
				},
				{
					Name: "Buggs Bunny",
				},
			},
			expectedOutput: []*User{
				{
					Name: "Buggs Bunny",
				},
				{
					Name: "Conor Wilson",
				},
				{
					Name: "Daffy Duck",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call function under test
			alphabetiseUsers(tt.users)

			// Confirm the output is as expected
			if len(tt.users) != len(tt.expectedOutput) {
				t.Errorf("alphabetiseUsers() returned a slice of unexpected length:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.users)
			}
			for i := 0; i < len(tt.expectedOutput); i++ {
				errMsg := fmt.Sprintf("alphabetiseUsers() returned an unexpected output:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.users)
				require.Equal(t, tt.users[i], tt.expectedOutput[i], errMsg)
			}
		})
	}
}

func TestContains(t *testing.T) {

	tests := []struct {
		name           string
		query          string
		usernames      []string
		expectedOutput bool
	}{
		{
			name: "does-contain",
			usernames: []string{
				"conor-wilson",
				"Daffy-Duck",
				"buggsBunny123",
			},
			query:          "Daffy-Duck",
			expectedOutput: true,
		},
		{
			name: "does-not-contain",
			usernames: []string{
				"conor-wilson",
				"buggsBunny123",
			},
			query:          "Daffy-Duck",
			expectedOutput: false,
		},
		{
			name:           "nil-usernames",
			usernames:      nil,
			query:          "Daffy-Duck",
			expectedOutput: false,
		},
	}
	// t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call function under test
			output := contains(tt.usernames, tt.query)

			// Confirm output is as expected
			if output != tt.expectedOutput {
				t.Errorf("contains() returned an unexpected output:\nWant: %v\nGot:  %v", tt.expectedOutput, output)
			}
		})
	}
}
