package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

type TreeDefinition struct {
	BaseName     string
	Descriptions []Description
}

type Description struct {
	Name   string
	Fields []ExpressionDescriptionField
}

type ExpressionDescriptionField struct {
	Name     string
	TypeName string
}

const visitorTemplate = `
type Visitor[T any] interface {
	{{ range .Descriptions }}
	Visit{{ .Name }}{{ $.BaseName }}(*{{ .Name }}{{ $.BaseName }}) T
  {{ end }}
}`

const expressionTemplate = `
type {{ .BaseName }} interface {
	getType() {{ .BaseName }}Type
}

type {{ $.BaseName }}Type int

const (
		{{ range $i,$v := .Descriptions }}
		{{ $v.Name | ToUpper }} {{ if eq $i 0}} = iota {{ end }}
		{{ end }}
)


{{ range .Descriptions }}
type {{ .Name }}{{ $.BaseName }} struct {
  {{ range .Fields }}
    {{ .Name }} {{ .TypeName }}
  {{ end }}
}

func (exp *{{ .Name }}{{ $.BaseName }}) getType() {{ $.BaseName }}Type {
	return {{ .Name | ToUpper }}
}

func (exp *{{ .Name }}{{ $.BaseName }}) Accept(visitor Visitor[any]) any {
	return {{ .Name }}{{ $.BaseName }}Accept(exp, visitor)
}

func {{ .Name }}{{ $.BaseName }}Accept[T any](expression *{{ .Name }}{{ $.BaseName }}, visitor Visitor[T]) T {
	return visitor.Visit{{ .Name }}{{ $.BaseName }}(expression)
}
{{ end }}


func expressionAccept[T any](e Expression, visitor Visitor[T]) T {
	switch e.getType() {
		{{ range $i,$v := .Descriptions }}
		case {{ $v.Name | ToUpper }}: 
			return {{ $v.Name }}{{ $.BaseName }}Accept(e.(*{{ $v.Name }}{{ $.BaseName }}), visitor)
		{{ end }}
	}
	return *new(T)
}
`

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
}

func DefineAST(outputDir string, treeDefinition TreeDefinition) error {
	visitorFile, err := os.Create(outputDir + "/visitor.go")

	if err != nil {
		return err
	}

	_, err = visitorFile.Write([]byte("package main\n"))

	if err != nil {
		return err
	}

	visitorTemplate, err := template.New("visitor").Parse(visitorTemplate)

	if err != nil {
		return err
	}

	expressionFile, err := os.Create(outputDir + "/expression.go")

	if err != nil {
		return err
	}

	_, err = expressionFile.Write([]byte("package main\n"))

	if err != nil {
		return err
	}

	expressionTmpl, err := template.New("expression").Funcs(funcMap).Parse(expressionTemplate)

	if err != nil {
		return err
	}

	err = visitorTemplate.Execute(visitorFile, treeDefinition)

	if err != nil {
		return err
	}

	err = expressionTmpl.Execute(expressionFile, treeDefinition)
	if err != nil {
		return err
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
					Name:     "Operator",
					TypeName: "Token",
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

	treeDefinition := TreeDefinition{
		BaseName:     "Expression",
		Descriptions: descriptions,
	}

	err := DefineAST(currentDir+"/..", treeDefinition)

	fmt.Println(err)
}
