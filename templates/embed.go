package templates

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed */*.template
var templateFs embed.FS

// Template is an interface for HTML templates.
type Template interface {
	// TemplateFilename description of the Go function.
	TemplateFilename() string
}

// Parser is a struct that represents a parser for HTML templates.
type Parser[T Template] struct {
	Template T
}

// NewParser creates a new Parser instance with the provided data.
//
// Parameters:
// - data: The data to be used for the Parser instance.
//
// Returns:
// - A pointer to the newly created Parser instance.
func NewParser[T Template](data T) *Parser[T] {
	return &Parser[T]{
		Template: data,
	}
}

// Parse is a Go function that parses the template and returns a string and an error.
func (ta *Parser[T]) Parse() (string, error) {
	return parseHtmlTemplate(ta.Template.TemplateFilename(), ta.Template)
}

// parseHtmlTemplate parses the HTML template with the given filename and data and returns the parsed template as a string or an error.
//
// filename string, data interface{}.
// string, error.
func parseHtmlTemplate(filename string, data interface{}) (string, error) {
	tmpl, err := template.ParseFS(templateFs, filename)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
