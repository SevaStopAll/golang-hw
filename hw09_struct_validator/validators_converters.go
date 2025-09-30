package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type (
	converter    = func(str string) (interface{}, error)
	validator    = func(f string, v interface{}, s ValidationStep) *ValidationError
	validatorSet map[string]validator
)

var (
	int64Type, strType = reflect.TypeOf((int64)(0)), reflect.TypeOf("")
	types              = []reflect.Type{int64Type, strType}
	converters         = map[reflect.Type]converter{
		strType: func(str string) (interface{}, error) {
			return str, nil
		},
		int64Type: func(str string) (interface{}, error) {
			return strconv.ParseInt(str, 0, 64)
		},
	}
	validators = map[reflect.Type]validatorSet{
		int64Type: {
			"max": validateIntMax,
			"min": validateIntMin,
			"in":  validateIntIn,
		},
		strType: {
			"len":    validateStrLen,
			"regexp": validateStrRegexp,
			"in":     validateStrIn,
		},
	}
)

func validateIntMax(f string, v interface{}, s ValidationStep) *ValidationError {
	var e error
	if maxInt, err := s.value(int64Type); err != nil {
		e = fmt.Errorf(
			"%w: field %s, Param %s isn't suitable for validator %s",
			ErrValidationFormat, f, s.Param, s.Op,
		)
	} else if val := v.(int64); val > maxInt.(int64) {
		e = ErrValidationMax
	}

	if e != nil {
		return &ValidationError{Field: f, Step: s, Val: v, Err: e}
	}

	return nil
}

func validateIntMin(f string, v interface{}, s ValidationStep) *ValidationError {
	var e error
	if minInt, err := s.value(int64Type); err != nil {
		e = fmt.Errorf(
			"%w: field %s, Param %s isn't suitable for validator %s",
			ErrValidationFormat, f, s.Param, s.Op,
		)
	} else if val := v.(int64); val < minInt.(int64) {
		e = ErrValidationMin
	}

	if e != nil {
		return &ValidationError{Field: f, Step: s, Val: v, Err: e}
	}
	return nil
}

func validateIntIn(f string, v interface{}, s ValidationStep) *ValidationError {
	return validateIn(f, v, s, int64Type)
}

func validateStrLen(f string, v interface{}, s ValidationStep) *ValidationError {
	var e error
	if l, err := s.value(int64Type); err != nil {
		e = fmt.Errorf(
			"%w: field %s, Param %s isn't suitable for validator %s",
			ErrValidationFormat, f, s.Param, s.Op,
		)
	} else if lv := int64(len(v.(string))); lv != l.(int64) {
		e = ErrValidationLen
	}

	if e != nil {
		return &ValidationError{Field: f, Step: s, Val: v, Err: e}
	}
	return nil
}

func validateStrRegexp(f string, v interface{}, s ValidationStep) *ValidationError {
	var e error
	re, err := regexp.Compile(s.Param)
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, Param %s isn't suitable for validator %s",
			ErrValidationFormat, f, s.Param, s.Op,
		)
	} else if !re.MatchString(v.(string)) {
		e = ErrValidationRegexp
	}

	if e != nil {
		return &ValidationError{Field: f, Step: s, Val: v, Err: e}
	}
	return nil
}

func validateStrIn(f string, v interface{}, s ValidationStep) *ValidationError {
	return validateIn(f, v, s, strType)
}

func validateIn(f string, v interface{}, s ValidationStep, t reflect.Type) *ValidationError {
	var e error
	set, err := s.slice(t)
	if err != nil {
		e = fmt.Errorf(
			"%w: field %s, Param %s isn't suitable for validator %s",
			ErrValidationFormat, f, s.Param, s.Op,
		)
	} else {
		var found bool
		for _, i := range set {
			if found = i == v; found {
				break
			}
		}
		if !found {
			e = ErrValidationIn
		}
	}

	if e != nil {
		return &ValidationError{Field: f, Step: s, Val: v, Err: e}
	}
	return nil
}
