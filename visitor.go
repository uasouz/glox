package main

type Visitor[T any] interface {
  
	VisitBinaryExpression(*BinaryExpression) T
  
	VisitGroupingExpression(*GroupingExpression) T
  
	VisitLiteralExpression(*LiteralExpression) T
  
	VisitUnaryExpression(*UnaryExpression) T
  
}