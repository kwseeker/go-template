package reflect

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

/*
借助标签实现的字段校验器
*/

const tagValidate = "validate"

// Regular expression to validate email address.
var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

// Validator Generic data validator
type Validator interface {
	//Validate method performs validation and returns results and optional error.
	Validate(interface{}) (bool, error)
}

// DefaultValidator does not perform any validations
type DefaultValidator struct {
}

func (v DefaultValidator) Validate(interface{}) (bool, error) {
	return true, nil
}

type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(int)

	if num < v.Min {
		return false, fmt.Errorf("should be greater than %v", v.Min)
	}

	if v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("should be less than %v", v.Max)
	}

	return true, nil
}

// StringValidator validates string presence and/or its length
type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	l := len(val.(string))

	if l == 0 {
		return false, fmt.Errorf("cannot be blank")
	}

	if l < v.Min {
		return false, fmt.Errorf("should be at least %v chars long", v.Min)
	}

	if v.Max >= v.Min && l > v.Max {
		return false, fmt.Errorf("should be less than %v chars long", v.Max)
	}

	return true, nil
}

type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}) (bool, error) {
	if !mailRe.MatchString(val.(string)) {
		return false, fmt.Errorf("is not a valid email address")
	}

	return true, nil
}

// Returns validator struct corresponding to validation type
func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		_, err := fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		if err != nil {
			return nil
		}
		return validator
	case "string":
		validator := StringValidator{}
		_, err := fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		if err != nil {
			return nil
		}
		return validator
	case "email":
		return EmailValidator{}
	}

	return DefaultValidator{}
}

// Performs actual data validation using validator definitions on the struct
func validateStruct(s interface{}) []error {
	//errs := []error{}
	var errs []error

	//ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		//Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagValidate)

		//Skip if tag is not defined or ignored
		if tag == "" || tag == "-" {
			continue
		}

		//Get a validator that corresponds to a tag
		validator := getValidatorFromTag(tag)

		//Perform validation
		valid, err := validator.Validate(v.Field(i).Interface())

		//Append error to results
		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}

	return errs
}
