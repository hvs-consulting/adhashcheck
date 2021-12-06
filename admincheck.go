package main

import "strings"

// Finds users whose name is substring of other user names
func findSubstringUsers(usernames []string) [][]string {
	reused := make([][]string, 0)

	allUsersLowercase := make([]string, len(usernames))
	for pos, username := range usernames {
		allUsersLowercase[pos] = strings.ToLower(username)
	}
	for _, username := range allUsersLowercase {
		reuse := []string{username}
		for _, toCompare := range allUsersLowercase {
			if toCompare != username && strings.Contains(toCompare, username) {
				reuse = append(reuse, toCompare)
			}
		}
		if len(reuse) > 1 {
			reused = append(reused, reuse)
		}
	}

	return reused
}
