package main

// NewField creates fields
func NewField(args []interface{}) *Field {
	return &Field{
		Name: args[0].(string),
		Type: args[1].(string),
	}
}

// Field is a field
type Field struct {
	Name string
	Type string
}
