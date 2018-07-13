package main

import "net/url"

// This is bad - do it right
func isValidURL(toTest string) bool {
	u, err := url.Parse(toTest)
	if err != nil {
		return false
	}

	if u.Scheme == "" {
		return false
	}

	return true
}
