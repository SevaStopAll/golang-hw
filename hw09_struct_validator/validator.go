package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrValueIsNotAStruct = errors.New("value is not a struct")
	ErrValidationStatus  = errors.New("validation completed")

	ErrValidationFormat          = errors.New("invalid format")
	ErrValidationTag             = fmt.Errorf("%w tag", ErrValidationFormat)
	ErrValidationUnsupportedRule = fmt.Errorf("%w rule", ErrValidationFormat)
	ErrValidationUnsupportedType = fmt.Errorf("%w type", ErrValidationFormat)

	ErrValidationFailed = errors.New("validation failed")
	ErrValidationMin    = fmt.Errorf("%w for validator Min", ErrValidationFailed)
	ErrValidationMax    = fmt.Errorf("%w for validator Max", ErrValidationFailed)
	ErrValidationIn     = fmt.Errorf("%w for validator In", ErrValidationFailed)
	ErrValidationLen    = fmt.Errorf("%w for validator Len", ErrValidationFailed)
	ErrValidationRegexp = fmt.Errorf("%w for validator Regexp", ErrValidationFailed)
)

type ValidationStep struct {
	Op    string
	Param string
}

func (vs ValidationStep) value(t reflect.Type) (interface{}, error) {
	return converters[t](vs.Param)
}

func (vs ValidationStep) slice(t reflect.Type) ([]interface{}, error) {
	res := make([]interface{}, 0, 5)
	for _, str := range strings.Split(vs.Param, ",") {
		if elem, err := converters[t](str); err == nil {
			res = append(res, elem)
		} else {
			return nil, err
		}
	}
	return res, nil
}

func newStep(str string) *ValidationStep {
	if parts := strings.SplitN(str, ":", 2); len(parts) == 2 {
		return &ValidationStep{
			Op:    parts[0],
			Param: parts[1],
		}
	}

	return nil
}

func parseTag(f string, str string) ([]ValidationStep, error) {
	if str == "" {
		return []ValidationStep{}, nil
	}
	stepLines := strings.Split(str, "|")
	steps := make([]ValidationStep, 0, len(stepLines))
	for _, v := range stepLines {
		if s := newStep(v); s != nil {
			steps = append(steps, *s)
			continue
		}
		return nil, fmt.Errorf("field '%s': '%s' is %w", f, str, ErrValidationTag)
	}

	return steps, nil
}

type ValidationError struct {
	Field string
	Step  ValidationStep
	Val   interface{}
	Err   error
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

func (v ValidationError) Error() string {
	if !errors.Is(v.Err, ErrValidationFailed) {
		return v.Err.Error()
	}

	message := fmt.Sprintf("%v - %s %s expected", v.Err, v.Step.Op, v.Step.Param)

	return fmt.Sprintf("Field '%s': value:'%v' message: %s", v.Field, v.Val, message)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	for _, ve := range v {
		builder.WriteString(ve.Error())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (v ValidationErrors) Unwrap() error {
	return ErrValidationStatus
}

func validateValue(field string, value reflect.Value, steps []ValidationStep) ValidationErrors {
	for _, t := range types {
		if value.Type().ConvertibleTo(t) {
			return callValidator(field, value.Convert(t).Interface(), steps, validators[t])
		}
	}
	if value.Type().Kind() == reflect.Slice {
		return validateSlice(field, value, steps)
	}
	return ValidationErrors{
		ValidationError{
			Field: field,
			Err:   fmt.Errorf("fgitield '%s': %w %s", field, ErrValidationUnsupportedType, value.Type()),
		},
	}
}

func validateSlice(field string, value reflect.Value, steps []ValidationStep) ValidationErrors {
	ev := make(ValidationErrors, 0, 3)
	for i := 0; i < value.Len(); i++ {
		if valErr := validateValue(field, value.Index(i), steps); valErr != nil {
			ev = append(ev, valErr...)
		}
	}
	return ev
}

func callValidator(f string, v interface{}, steps []ValidationStep, vSet validatorSet) ValidationErrors {
	ev := make(ValidationErrors, 0, 3)
	for _, s := range steps {
		if vf := vSet[s.Op]; vf == nil {
			ev = append(ev, ValidationError{
				Field: f, Step: s, Val: v,
				Err: fmt.Errorf("field '%s': %w '%s' for field '%s'", f, ErrValidationUnsupportedRule, s.Op, f),
			})
			continue
		} else if valErr := vf(f, v, s); valErr != nil {
			ev = append(ev, *valErr)
		}
	}

	return ev
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("%v %w", v, ErrValueIsNotAStruct)
	}
	ve := make(ValidationErrors, 0, 10)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("validate")
		if tag == "" {
			continue
		}
		if steps, e := parseTag(f.Name, tag); e == nil {
			if valErr := validateValue(f.Name, val.Field(i), steps); valErr != nil {
				ve = append(ve, valErr...)
			}
		} else {
			ve = append(ve, ValidationError{Field: f.Name, Err: e, Val: val.Field(i).Interface()})
		}
	}
	return ve
}
