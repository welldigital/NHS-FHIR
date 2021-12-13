package client

import "fmt"

// RateLimitError contains information relating to this type of error
type RateLimitError struct {
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintln("You have exceeeded the rate limit for this API.")
}
