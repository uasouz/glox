package main

import "fmt"

type AstPrinter struct {
}

func (p *AstPrinter) Print(expression Expression) string {
	return expressionAccept(expression, p)
}

func (p *AstPrinter) VisitBinaryExpression(expression *BinaryExpression) (_ string) {
	return p.parenthesize(expression.Operator.lexeme, expression.Left, expression.Right)
}

func (p *AstPrinter) VisitGroupingExpression(expression *GroupingExpression) (_ string) {
	return p.parenthesize("group", expression.Expression)
}

func (p *AstPrinter) VisitLiteralExpression(expression *LiteralExpression) (_ string) {
	if expression.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expression.Value)
}

func (p *AstPrinter) VisitUnaryExpression(expression *UnaryExpression) (_ string) {
	return p.parenthesize(expression.Operator.lexeme, expression.Right)
}

func (p *AstPrinter) parenthesize(lexeme string, expressions ...Expression) string {

	result := "( " + lexeme

	for _, expression := range expressions {
		result += " "
		result += expressionAccept(expression, p)
	}

	result += ")"
	return result
}

