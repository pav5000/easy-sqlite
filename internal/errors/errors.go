// errors provides a couple of custom useful error handling functions
// parts of the code are taken from github.com/pkg/errors
// the pkg/errors repo is archived and deprecated so I didn't want to import it
package errors

import "errors"

var New = errors.New

// Wpr annotates err with a new message.
// If err is nil, WithMessage returns nil.
func Wrp(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		err: err,
		msg: message,
	}
}

type withMessage struct {
	err error
	msg string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.err.Error() }

func (w *withMessage) Unwrap() error { return w.err }
