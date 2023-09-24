package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name             string  `json:"name"`
	Login            string  `json:"login"`
	Company          string  `json:"company"`
	NumFollowers     int     `json:"followers"`
	NumRepos         int     `json:"public_repos"`
	AvgRepoFollowers float32 `json:"avg_public_repo_followers"`
}

type Users []*User

func main() {
	router := gin.Default()
	router.GET("/retrieveUsers", retrieveUsers)
	router.Run("localhost:8080")
}

// retrieveUsers retrieves the basic user information of the specified GitHub users. This
// function encapsulates the functionality of this API's retrieveUsers endpoint.
func retrieveUsers(c *gin.Context) {

	// Get the array of usernames from the query
	usernames := c.QueryArray("users")
	usernames = deduplicate(usernames)

	// For each username...
	users := Users{}
	for _, username := range usernames {

		// ...obtain the raw user data from GitHub's API...
		resp, err := http.Get("https://api.github.com/users/" + username)
		if err != nil {
			log.Printf("error accesssing GitHub's API: %v\n", err)
			return
		}

		// ...unmarshal the resulting response into a User struct...
		newUser, success, err := unmarshalUserFromGitHubResponse(resp)
		if err != nil {
			log.Printf("error unmarshalling User from GitHub response: %v\n", err)
			return
		} else if !success {
			log.Printf("[WARNING] Request for user information with username '%v' was unsuccessful. Status:, %v\n", username, resp.Status)
			continue
		}
		newUser.computeAvgRepoFollowers()

		// ...and append the new User to the User slice.
		users = append(users, newUser)
	}

	// Tidy up the result and push the information to the API output.
	users.alphabetise()
	c.IndentedJSON(http.StatusOK, users)
}

// unmarshalUserFromGitHubResponse marshals a User struct from the body of the
// http.Response struct returned by GitHub's API.
func unmarshalUserFromGitHubResponse(resp *http.Response) (*User, bool, error) {

	// Check to confirm that the request was successful (if not, we simply log a warning
	// and move on)
	if resp.StatusCode != 200 {
		return nil, false, nil
	}

	// Convert response to a slice of JSON bytes.
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response from GitHub: %v", err)
	}

	// Marshal the JSON to our User struct type
	newUser := &User{}
	if err = json.Unmarshal([]byte(respData), newUser); err != nil {
		return nil, false, fmt.Errorf("error unmarshalling from JSON: %v\n", err)
	}
	return newUser, true, nil
}

// deduplicate returns a slice of string usernames identical to the provided slice of
// string usernames, but with duplicates removed.
func deduplicate(usernames []string) []string {
	output := []string{}
	for _, username := range usernames {
		if !contains(output, username) {
			output = append(output, username)
		}
	}
	return output
}

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

// contains indicates whether or not the provided slice of string usernames contains the
// query string username.
func contains(usernames []string, query string) bool {
	for _, existingString := range usernames {
		if existingString == query {
			return true
		}
	}
	return false
}
