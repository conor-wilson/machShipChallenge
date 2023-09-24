package main

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
