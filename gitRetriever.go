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
	Name             string `json:"name"`
	Login            string `jason:"login"`
	Company          string `json:"company"`
	NumFollowers     int    `json:"followers"`
	NumRepos         int    `json:"public_repos"`
	AvgRepoFollowers float32
}

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
	users := []User{}
	for _, username := range usernames {

		// ...obtain the raw user data from GitHub...
		userData, success, err := getUserDataFromGitHub(username)
		if err != nil {
			fmt.Println(err)
			return
		} else if !success {
			continue
		}

		// ...and marshal the data into our User struct type.
		newUser := User{}
		if err = json.Unmarshal([]byte(userData), &newUser); err != nil {
			fmt.Printf("Error unmarshalling from JSON: %v\n", err)
			return
		}
		users = append(users, newUser)
	}

	// Tidy up the result and push the information to the API output.
	alphabetiseUsers(users)
	c.IndentedJSON(http.StatusOK, users)
}

// getUserDataFromGitHub obtains the raw user data of the provided user from GitHub and
// outputs it as a slice of JSON bytes.
func getUserDataFromGitHub(username string) ([]byte, bool, error) {

	// Get the data response from GitHub's API
	resp, err := http.Get("https://api.github.com/users/" + username)
	if err != nil {
		return nil, false, fmt.Errorf("error accesssing GitHub's API: %v", err)
	}

	// Check to confirm that the request was successful (if not, we simply log a warning
	// and move on)
	if resp.StatusCode != 200 {
		log.Printf("[WARNING] Request for user information with username '%v' was unsuccessful. Status:, %v", username, resp.Status)
		return nil, false, nil
	}

	// Convert response to a slice of JSON bytes.
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("error reading response from GitHub: %v", err)
	}

	return respData, true, nil
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

// alphabetiseUsers sorts the provided slice of Users alphabetically by name.
func alphabetiseUsers(users []User) {
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
