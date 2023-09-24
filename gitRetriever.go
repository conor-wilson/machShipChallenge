package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
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
	users := []user{}

	// For each username...
	for _, username := range usernames {

		// ...obtain the raw user data from GitHub...
		userData, err := getUserDataFromGitHub(username)
		if err != nil {
			fmt.Println(err)
			return
		}

		// ...and marshal the data into our specified user struct type.
		newUser := user{}
		if err = json.Unmarshal([]byte(userData), &newUser); err != nil {
			fmt.Printf("Error unmarshalling from JSON: %v\n", err)
			return
		}
		users = append(users, newUser)
	}

	// Push the information to the API output.
	c.IndentedJSON(http.StatusOK, users)
}

// getUserDataFromGitHub obtains the raw user data of the provided user from GitHub and
// outputs it as a slice of JSON bytes.
func getUserDataFromGitHub(username string) ([]byte, error) {

	// Get the data response from GitHub's API
	resp, err := http.Get("https://api.github.com/users/" + username)
	if err != nil {
		return nil, fmt.Errorf("error accesssing GitHub's API: %v", err)
	}

	// Convert response to a slice of JSON bytes.
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from GitHub: %v", err)
	}
	return respData, nil
}
