package errors

import (
	"bytes"
	"fmt"
)

// Unwrap recursively finds non-wrap error and returns the first one encountered
// (e.g NotFound, Conflict etc...) or returns the error untouched if it is not wrapped
func Unwrap(err error) error {
	u, ok := err.(interface{
		Unwrap() error
	})
	if !ok {
		return err
	}
	return Unwrap(u.Unwrap())
}

// Wrap adds context to already wrapped error or wraps it if it is not.
// If error is nil, nil is returned. This makes it easy to write expressions
// like this: return u, error.Wrap(tx.Commit(), op, "...")
func Wrap(err error, op string, message string) error {
	if err == nil {
		return nil
	}
	return &wrap{Err: err, Op:op, privMsg:message}
}

type wrap struct {
	Err     error
	Op      string
	privMsg string
}

func (c *wrap) Error() string {
	var buf bytes.Buffer

	// Print the current operation in our stack, if any.
	if c.Op != "" {
		_, _ = fmt.Fprintf(&buf, "[%s]: ", c.Op)
	} else {
		_, _ = fmt.Fprint(&buf, "_: ")
	}

	// Print the current additional context in our stack, if any.
	if c.privMsg != "" {
		_, _ = fmt.Fprintf(&buf, "[%s] >> ", c.privMsg)
	} else {
		_, _ = fmt.Fprint(&buf, "_ >> ")
	}

	// If wrapping an error, print its Error() message. Otherwise print the error code & message.
	if c.Err != nil {
		buf.WriteString(c.Err.Error())
	} else {
		_, _ = fmt.Fprintf(&buf, "<Generic error> ")
		buf.WriteString(c.privMsg)
	}

	return buf.String()
}

func(c *wrap) Unwrap() error {
	return c.Err
}
