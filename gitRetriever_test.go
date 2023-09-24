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
