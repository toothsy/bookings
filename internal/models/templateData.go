package models

import "github/toothsy/bookings/internal/forms"

// TemplateData holds data that can be passed to tmpl files
type TemplateData struct {
	StringMap map[string]string
	NumberMap map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	Warning   string
	Info      string
	Error     string
	CSRFToken string
	Form      *forms.Form
}
