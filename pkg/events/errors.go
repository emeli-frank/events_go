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

type Conflict struct {
	Err     error
}

// Error outputs stack info that should not be shown to client.
func (e *Conflict) Error() string {
	return e.Err.Error()
}

func (e *Conflict) Cause() error {
	return e.Err
}

type Unauthorized struct {
	Err     error
}

// Error outputs stack info that should not be shown to client.
func (e *Unauthorized) Error() string {
	return e.Err.Error()
}

func (e *Unauthorized) Cause() error {
	return e.Err
}

type Invalid struct {
	Err     error
}

// Error outputs stack info that should not be shown to client.
func (e *Invalid) Error() string {
	return e.Err.Error()
}

func (e *Invalid) Cause() error {
	return e.Err
}
