package main

import "sort"

// User encapsulates the basic information of a GitHub user.
type User struct {
	Name         string `json:"name"`
	Login        string `json:"login"`
	Company      string `json:"company"`
	NumFollowers int    `json:"followers"`
	NumRepos     int    `json:"public_repos"`

	// The average number of followers the user has per public repository.
	AvgRepoFollowers float32 `json:"avg_public_repo_followers"`
}

// Users encapsulates a slice of Users.
type Users []*User

// computeAvgRepoFollowers calculates the average followers per public repo of the
// provided User and update's the User's field accordingly.
func (user *User) computeAvgRepoFollowers() {
	if user.NumRepos == 0 {
		user.AvgRepoFollowers = 0
		return
	}
	user.AvgRepoFollowers = float32(user.NumFollowers) / float32(user.NumRepos)
}

// alphabetise sorts the provided slice of Users alphabetically by name.
func (users Users) alphabetise() {
	sort.Slice(users, func(i, j int) bool {
		return users[i].Name[0] < users[j].Name[0]
	})
}
