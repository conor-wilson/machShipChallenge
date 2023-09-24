package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// main routes the functionality of the API and exposes it on port 8080.
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

	users := Users{}
	for _, username := range usernames {

		// Obtain the raw user data from GitHub's API
		resp, err := http.Get("https://api.github.com/users/" + username)
		if err != nil {
			log.Printf("error accesssing GitHub's API: %v\n", err)
			return
		}

		// Unmarshal GitHub's response into a User struct
		newUser, success, err := unmarshalUserFromGitHubResponse(resp)
		if err != nil {
			log.Printf("error unmarshalling User from GitHub response: %v\n", err)
			return
		} else if !success {
			log.Printf("[WARNING] Request for user information with username '%v' was unsuccessful. Status:, %v\n", username, resp.Status)
			continue
		}

		// Tidy the fields and append the new User to the slice of Users.
		newUser.computeAvgRepoFollowers()
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
