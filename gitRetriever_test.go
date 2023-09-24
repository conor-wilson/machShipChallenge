package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalUserFromGitHubResponse(t *testing.T) {

	testUser := &User{
		Name:             "Conor Wilson",
		Login:            "conor-wilson",
		Company:          "MachShip",
		NumFollowers:     23,
		NumRepos:         4,
		AvgRepoFollowers: float32(23) / float32(4),
	}
	testUserJSON, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("error marshalling test User to JSON: %v", err)
	}

	tests := []struct {
		name            string
		resp            *http.Response
		expectedUser    *User
		expectedSuccess bool
		errorExists     bool
	}{
		{
			name: "happy-path",
			resp: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(string(testUserJSON))),
			},
			expectedUser:    testUser,
			expectedSuccess: true,
			errorExists:     false,
		},
		{
			name: "happy-path-non-existant-username",
			resp: &http.Response{
				StatusCode: 404,
				Status:     "404 Not Found",
				Body:       nil,
			},
			expectedUser:    nil,
			expectedSuccess: false,
			errorExists:     false,
		},
		{
			name: "err-unmarshalling-json",
			resp: &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("")),
			},
			expectedUser:    nil,
			expectedSuccess: false,
			errorExists:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call function under test
			user, success, err := unmarshalUserFromGitHubResponse(tt.resp)

			// Confirm the output is as expected
			errMsg := fmt.Sprintf("unmarshalUserFromGitHubResponse() returned an unexpected User struct:\nWant: %v\nGot:  %v", tt.expectedUser, user)
			require.Equal(t, user, tt.expectedUser, errMsg)
			if success != tt.expectedSuccess {
				t.Errorf("unmarshalUserFromGitHubResponse() returned an unexpected bool var:\nWant: %v\nGot:  %v", tt.expectedSuccess, success)
			} else if err == nil && tt.errorExists {
				t.Errorf("unmarshalUserFromGitHubResponse() did not return an error when it should have.")
			} else if err != nil && !tt.errorExists {
				t.Errorf("unmarshalUserFromGitHubResponse() returned an error when it shouldn't have: %v", err)
			}
		})
	}
}

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

func TestComputeAvgRepoFollowers(t *testing.T) {

	tests := []struct {
		name           string
		user           *User
		expectedOutput *User
	}{
		{
			name: "happy-path",
			user: &User{
				NumFollowers: 23,
				NumRepos:     4,
			},
			expectedOutput: &User{
				NumFollowers:     23,
				NumRepos:         4,
				AvgRepoFollowers: float32(23) / float32(4),
			},
		},
		{
			name: "happy-path-0-followers",
			user: &User{
				NumFollowers: 0,
				NumRepos:     4,
			},
			expectedOutput: &User{
				NumFollowers:     0,
				NumRepos:         4,
				AvgRepoFollowers: 0,
			},
		},
		{
			name: "happy-path-0-repos",
			user: &User{
				NumFollowers: 23,
				NumRepos:     0,
			},
			expectedOutput: &User{
				NumFollowers:     23,
				NumRepos:         0,
				AvgRepoFollowers: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call function under test
			computeRepoAvgFollowers(tt.user)

			// Confirm the output is as expected
			errMsg := fmt.Sprintf("computeAverageFollowers() did not compute the average repository followers correctly:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.user)
			require.Equal(t, tt.user, tt.expectedOutput, errMsg)
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
