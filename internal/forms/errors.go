package forms

// the course does not do pointer receiver
type errors map[string][]string

// Add adds the error message to related field
func (e *errors) Add(field, message string) {

	(*e)[field] = append((*e)[field], message)

}

// Get returns the error associated with the field
func (e *errors) Get(field string) string {
	hasField := (*e)[field]
	if len(hasField) == 0 {
		return ""
	}
	return hasField[0]
}
