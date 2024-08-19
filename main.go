package main

import (
	"bufio"
	"fmt"
	"os"
)

var hadError = false

func main() {
	testExpression := BinaryExpression{
		Left: &UnaryExpression{
			Token{
				Type:    MINUS,
				lexeme:  "-",
				literal: nil,
				line:    1,
			},
			&LiteralExpression{
				Value: 123,
			},
		},
		Operator: Token{
			Type: STAR,
			lexeme: "*",
			literal: nil,
			line: 1,
		},
		Right: &GroupingExpression{
			&LiteralExpression{
				Value: 45.67,
			},
		},
	}

	printer := AstPrinter{}

	fmt.Println(printer.Print(&testExpression))

	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage glox [script]")
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}

}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = run(string(bytes))
	if err != nil {
		os.Exit(65)
	}
	return err
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		run(line)
		hadError = false

		fmt.Print("> ")
	}

}

func run(source string) error {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

func Error(lineNumber int, message string) {
	report(lineNumber, "", message)
}

func report(lineNumber int, where string, message string) {
	fmt.Println(fmt.Sprintf("[line %d] Error %s: %s", lineNumber, where, message))
}
