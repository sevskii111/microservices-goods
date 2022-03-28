package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Form struct {
	url.Values
	Errors FormErrors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		FormErrors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)

	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too shot (minimum is %d characters)", d))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) MinValue(field string, minValue int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		f.Errors.Add(field, "This field must be numeric")
		return
	}

	if intValue < minValue {
		f.Errors.Add(field, fmt.Sprintf("This field must be >= %d", minValue))
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Has(field string) bool {
	return len(f.Values[field]) > 0
}
