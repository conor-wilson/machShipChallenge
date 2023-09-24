package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserComputeAvgRepoFollowers(t *testing.T) {

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
			name: "happy-path-zero-followers",
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
			name: "happy-path-zero-repos",
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
			tt.user.computeAvgRepoFollowers()

			// Confirm the output is as expected
			errMsg := fmt.Sprintf("User.computeAvgRepoFollowers() did not compute the average repository followers correctly:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.user)
			require.Equal(t, tt.user, tt.expectedOutput, errMsg)
		})
	}
}

func TestUsersAlphabetise(t *testing.T) {

	tests := []struct {
		name           string
		users          Users
		expectedOutput Users
	}{
		{
			name: "happy-path",
			users: Users{
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
			expectedOutput: Users{
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
			tt.users.alphabetise()

			// Confirm the output is as expected
			if len(tt.users) != len(tt.expectedOutput) {
				t.Errorf("Users.alphabetise() returned a slice of unexpected length:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.users)
			}
			for i := 0; i < len(tt.expectedOutput); i++ {
				errMsg := fmt.Sprintf("Users.alphabetise() returned an unexpected output:\nWant: %v\nGot:  %v", tt.expectedOutput, tt.users)
				require.Equal(t, tt.users[i], tt.expectedOutput[i], errMsg)
			}
		})
	}
}
