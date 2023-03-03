package forms

import (
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

// Valid checks if there are any errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks the form data's equivalent of the field
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Get(field)
		if strings.TrimSpace(val) == "" {
			f.Errors.Add(field, "This field cant be blank")
		}
	}
}

// Has checks if the request contains the field in it
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x == ""
}
