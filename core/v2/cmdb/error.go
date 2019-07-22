package cmdb

import (
	"errors"
	"fmt"
	"runtime"

	"strings"
)

var (
	ErrArgumentsMissing        = errors.New("arguments missing")
	ErrFailedToCreateInfoCMDB  = errors.New("failed to create infocmdb object")
	ErrNoCredentials           = errors.New("must provide credentials")
	ErrNotImplemented          = errors.New("not implemented")
	ErrNoResult                = errors.New("query returned no result")
	ErrTooManyResults          = errors.New("query returned to many results, expected one")
	ErrWebserviceResponseNotOk = errors.New("webservice response was not ok")
	ErrLoginFailed             = errors.New("login failed")
)

// Errors contains all happened errors
type Errors []error


func FunctionError(msg string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fullMsg :=  fmt.Sprintf("[%s:%d]:%s", frame.Function, frame.Line , msg)

	return errors.New(fullMsg)
}

func (i *InfoCMDB) AddError(err error) error {
	if err != nil {
		errorMessage := FunctionError(err.Error())
		i.Logger.Error(errorMessage)

		errors := Errors(i.GetErrors())
		errors = errors.Add(errorMessage)
		if len(errors) > 1 {
			err = errors
		}

		i.Error = err
	}
	return err
}


// GetErrors
func (i *InfoCMDB) GetErrors() []error {
	if errs, ok := i.Error.(Errors); ok {
		return errs
	} else if i.Error != nil {
		return []error{i.Error}
	}
	return []error{}
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