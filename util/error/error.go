package error

import (
	"errors"
	"runtime"
	"strings"
)

// Errors contains all happened errors
type Errors []error

func FunctionError(msg string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fullMsg := frame.Function + ": " + msg

	return errors.New(fullMsg)
}




// Add adds an error to a given slice of errors
func (errs Errors) Add(newErrors ...error) Errors {
	for _, err := range newErrors {
		if err == nil {
			continue
		}

		if errors, ok := err.(Errors); ok {
			errs = errs.Add(errors...)
		} else {
			ok = true
			for _, e := range errs {
				if err == e {
					ok = false
				}
			}
			if ok {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// GetErrors gets all errors that have occurred and returns a slice of errors (Error type)
func (errs Errors) GetErrors() []error {
	return errs
}

// Error takes a slice of all errors that have occurred and returns it as a formatted string
func (errs Errors) Error() string {
	var errors = []string{}
	for _, e := range errs {
		errors = append(errors, e.Error())
	}
	return strings.Join(errors, "; ")
}