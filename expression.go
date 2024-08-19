package main

type Expression interface {
	getType() ExpressionType
}

type ExpressionType int

const (
	BINARY = iota

	GROUPING

	LITERAL

	UNARY
)

type BinaryExpression struct {
	Left Expression

	Right Expression

	Operator Token
}

func (exp *BinaryExpression) getType() ExpressionType {
	return BINARY
}

func (exp *BinaryExpression) Accept(visitor Visitor[any]) any {
	return BinaryExpressionAccept(exp, visitor)
}

func BinaryExpressionAccept[T any](expression *BinaryExpression, visitor Visitor[T]) T {
	return visitor.VisitBinaryExpression(expression)
}

type GroupingExpression struct {
	Expression Expression
}

func (exp *GroupingExpression) getType() ExpressionType {
	return GROUPING
}

func (exp *GroupingExpression) Accept(visitor Visitor[any]) any {
	return GroupingExpressionAccept(exp, visitor)
}

func GroupingExpressionAccept[T any](expression *GroupingExpression, visitor Visitor[T]) T {
	return visitor.VisitGroupingExpression(expression)
}

type LiteralExpression struct {
	Value any
}

func (exp *LiteralExpression) getType() ExpressionType {
	return LITERAL
}

func (exp *LiteralExpression) Accept(visitor Visitor[any]) any {
	return LiteralExpressionAccept(exp, visitor)
}

func LiteralExpressionAccept[T any](expression *LiteralExpression, visitor Visitor[T]) T {
	return visitor.VisitLiteralExpression(expression)
}

type UnaryExpression struct {
	Operator Token

	Right Expression
}

func (exp *UnaryExpression) getType() ExpressionType {
	return UNARY
}

func (exp *UnaryExpression) Accept(visitor Visitor[any]) any {
	return UnaryExpressionAccept(exp, visitor)
}

func UnaryExpressionAccept[T any](expression *UnaryExpression, visitor Visitor[T]) T {
	return visitor.VisitUnaryExpression(expression)
}

func expressionAccept[T any](e Expression, visitor Visitor[T]) T {
	switch e.getType() {

	case BINARY:
		return BinaryExpressionAccept(e.(*BinaryExpression), visitor)

	case GROUPING:
		return GroupingExpressionAccept(e.(*GroupingExpression), visitor)

	case LITERAL:
		return LiteralExpressionAccept(e.(*LiteralExpression), visitor)

	case UNARY:
		return UnaryExpressionAccept(e.(*UnaryExpression), visitor)

	}
	return *new(T)
}
