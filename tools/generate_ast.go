package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type Description struct {
	Name   string
	Fields []ExpressionDescriptionField
}

type ExpressionDescriptionField struct {
	Name     string
	TypeName string
}

func visitorTemplate(baseName string) string {
	return strings.ReplaceAll(`
type Visitor[T any] interface {
  {{ range . }}
	Visit{{ .Name }}||baseName||(*{{ .Name }}||baseName||) T
  {{ end }}
}`, "||baseName||", baseName)
}

func expressionInterfaceTemplate(baseName string) string {
	return strings.ReplaceAll(`
type ||baseName|| interface {
	getType() string
}`, "||baseName||", baseName)
}

func expressionTemplate(baseName string) string {
	return strings.ReplaceAll(`
type {{ .Name }}||baseName|| struct {
  {{ range .Fields }}
    {{ .Name }} {{ .TypeName }}
  {{ end }}
}

func (exp *{{ .Name }}||baseName||) getType() string {
	return "{{ .Name | ToUpper }}"
}

func (exp *{{ .Name }}||baseName||) Accept(visitor Visitor[any]) any {
	return {{ .Name }}||baseName||Accept(exp, visitor)
}

func {{ .Name }}||baseName||Accept[T any](expression *{{ .Name }}||baseName||, visitor Visitor[T]) T {
	return visitor.Visit{{ .Name }}||baseName||(expression)
}
`, "||baseName||", baseName)
}

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}

func DefineAST(outputDir, baseName string, descriptions []Description) error {
	visitorFile, err := os.Create(outputDir + "/visitor.go")

	if err != nil {
		return err
	}

	_, err = visitorFile.Write([]byte("package main\n"))

	if err != nil {
		return err
	}

	visitorTemplate, err := template.New("visitor").Parse(visitorTemplate(baseName))

	if err != nil {
		return err
	}

	expressionFile, err := os.Create(outputDir + "/expression.go")

	if err != nil {
		return err
	}

	_, err = expressionFile.Write([]byte("package main\n" + expressionInterfaceTemplate(baseName) + "\n"))

	if err != nil {
		return err
	}

	expressionTmpl, err := template.New("expression").Funcs(funcMap).Parse(expressionTemplate(baseName))

	if err != nil {
		return err
	}

	err = visitorTemplate.Execute(visitorFile, descriptions)

	if err != nil {
		return err
	}

	for _, description := range descriptions {
		err = expressionTmpl.Execute(expressionFile, description)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	currentDir, _ := os.Getwd()

	descriptions := []Description{
		{
			Name: "Binary",
			Fields: []ExpressionDescriptionField{
				{
					Name:     "Left",
					TypeName: "Expression",
				},
				{
					Name:     "Right",
					TypeName: "Expression",
				},
				{
					Name:     "Lexeme",
					TypeName: "string",
				},
			},
		},
		{
			Name: "Grouping",
			Fields: []ExpressionDescriptionField{
				{
					Name:     "Expression",
					TypeName: "Expression",
				},
			},
		},
		{
			Name: "Literal",
			Fields: []ExpressionDescriptionField{
				{
					Name:     "Value",
					TypeName: "any",
				},
			},
		},
		{
			Name: "Unary",
			Fields: []ExpressionDescriptionField{
				{
					Name:     "Operator",
					TypeName: "Token",
				},
				{
					Name:     "Right",
					TypeName: "Expression",
				},
			},
		},
	}

	err := DefineAST(currentDir+"/..", "Expression", descriptions)

	fmt.Println(err)
}
