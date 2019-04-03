package protocol

import "fmt"

type NotALeaderError struct {
}

func (e *NotALeaderError) Error() string {
	return fmt.Sprintf("Node is not an error")
}
