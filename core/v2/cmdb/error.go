package cmdb

import (
	util_error "github.com/infonova/infocmdb-lib-go/util/error"
)


func (i *InfoCMDB) AddError(err error) error {
	if err != nil {
		errorMessage := util_error.FunctionError(err.Error())
		i.Logger.Error(errorMessage)

		errors := util_error.Errors(i.GetErrors())
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
	if errs, ok := i.Error.(util_error.Errors); ok {
		return errs
	} else if i.Error != nil {
		return []error{i.Error}
	}
	return []error{}
}