package events

type NotFound struct {
	Err     error
}

// Error outputs stack info that should not be shown to client.
func (e *NotFound) Error() string {
	return e.Err.Error()
}

func (e *NotFound) Cause() error {
	return e.Err
}
