package parser

import (
	"aura/src/ast"
	l "aura/src/lexer"
	"fmt"
	"strconv"
)

// parse a boolean expression
func (p *Parser) parseBoolean() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	var value bool
	if p.currentToken.Token_type == l.TRUE {
		value = true
		return ast.NewBoolean(*p.currentToken, &value)
	}

	value = false
	return ast.NewBoolean(*p.currentToken, &value)
}

// parse a for expression
func (p *Parser) parseFor() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	forExpression := ast.NewFor(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> por i en rango(10))
		return nil
	}

	forExpression.Condition = p.parseRangeExpression()
	if !p.expepectedToken(l.RPAREN) {
		// syntax error -> por(i en range(10)
		return nil
	}
	if !p.expepectedToken(l.LBRACE) {
		// syntax error -> por(i en rango(10))
		//					body...
		//				   }
		return nil
	}

	forExpression.Body = p.parseBlock()
	return forExpression
}

// parse a function expression
func (p *Parser) parseFunction() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	function := ast.NewFunction(*p.currentToken, nil)
	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> funcion {}
		return nil
	}

	function.Parameters = p.parseFunctionParameters()
	if !p.expepectedToken(l.LBRACE) {
		// syntax error -> funcion() there is on body
		return nil
	}

	function.Body = p.parseBlock()
	return function
}

// parse a while expression
func (p *Parser) parseWhile() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	whileExpression := ast.NewWhile(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> mientras <condition>
		// missing the left paren
		return nil
	}

	p.advanceTokens()
	whileExpression.Condition = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		// syntax error -> mientras <condition>) {}
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		// syntax error -> mientras (<condition>) there is no body
		return nil
	}

	whileExpression.Body = p.parseBlock()
	return whileExpression
}

// parse a identifier expression
func (p *Parser) parseIdentifier() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return ast.NewIdentifier(*p.currentToken, p.currentToken.Literal)
}

// parse an if expresion
func (p *Parser) parseIf() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	ifExpression := ast.NewIf(*p.currentToken, nil, nil, nil)
	if !p.expepectedToken(l.LPAREN) {
		// syntax error. missing parents
		return nil
	}

	p.advanceTokens()
	ifExpression.Condition = p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		// syntax error. missing parents
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		// syntax error. missing body
		return nil
	}

	ifExpression.Consequence = p.parseBlock()

	p.checkPeekTokenIsNotNil()
	// if we have an else token that means there is an else expression
	if p.peekToken.Token_type == l.ELSE {
		p.advanceTokens()
		if !p.expepectedToken(l.LBRACE) {
			return nil
		}

		ifExpression.Alternative = p.parseBlock()
	}

	return ifExpression
}

// parse a integer expressions
func (p *Parser) parseInteger() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	integer := ast.NewInteger(*p.currentToken, nil)

	val, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		// the value is not a number. this is very weird to happend
		message := fmt.Sprintf("no se pudo parsear %s como entero", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	integer.Value = &val
	return integer
}

// parse a float expression
func (p *Parser) parseFloat() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	float := &ast.FloatExp{BaseNode: ast.BaseNode{Token: *p.currentToken}}
	val, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		// the value is not a float. this is very weird to happend
		message := fmt.Sprintf("no se pudo parsear %s como flotante", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	float.Value = val
	return float
}

// parse a group expression like (5 + 5) / 2
func (p *Parser) parseGroupExpression() ast.Expression {
	p.advanceTokens()
	expression := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		// syntax error: missing parenthessis
		return nil
	}

	return expression
}

// parse a prefix expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	prefixExpression := ast.NewPrefix(*p.currentToken, p.currentToken.Literal, nil)
	p.advanceTokens()
	prefixExpression.Rigth = p.parseExpression(PREFIX)
	return prefixExpression
}

// parse a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return ast.NewStringLiteral(*p.currentToken, p.currentToken.Literal)
}

// parse a array expression
func (p *Parser) ParseArray() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	arr := ast.NewArray(*p.currentToken, nil)
	if !p.expepectedToken(l.LBRACKET) {
		// syntax error -> lista 2,3,4,5
		return nil
	}

	arr.Values = p.ParseArrayValues()
	return arr
}

// parse a map expression
func (p *Parser) parseMap() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	mapExpress := ast.NewMapExpression(*p.currentToken, []*ast.KeyValue{})
	if !p.expepectedToken(l.LBRACE) {
		// syntax error: missing left brace
		return nil
	}

	p.advanceTokens()
	keyVal := p.parseKeyValues()
	if keyVal != nil {
		mapExpress.Body = append(mapExpress.Body, keyVal)
	}

	// we loop untile we dont have commas. this means we parse all the key value pairs.
	for p.peekToken.Token_type == l.COMMA {
		p.advanceTokens()
		p.advanceTokens()
		keyVal := p.parseKeyValues()
		if keyVal != nil {
			mapExpress.Body = append(mapExpress.Body, keyVal)
		}

	}
	p.advanceTokens()
	return mapExpress
}

func (p *Parser) parseClassMethod() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	classMethod := ast.NewClassMethodExp(*p.currentToken, nil, nil, nil)
	classMethod.Name = p.parseIdentifier().(*ast.Identifier)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	classMethod.Params = p.parseFunctionParameters()
	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	classMethod.Body = p.parseBlock()
	p.advanceTokens()
	return classMethod
}

func (p *Parser) parseClassField() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	field := ast.NewClassFieldExp(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.DOT) {
		return nil
	}

	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	field.Name = p.parseIdentifier().(*ast.Identifier)
	if p.peekToken.Token_type != l.ASSING {
		return field
	}

	p.advanceTokens()
	p.advanceTokens()
	field.Value = p.parseExpression(LOWEST)

	p.checkPeekTokenIsNotNil()
	if p.peekToken.Token_type == l.SEMICOLON {
		p.advanceTokens()
	}

	return field
}

func (p *Parser) parseClassCall() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	call := ast.NewClassCall(*p.currentToken, nil, nil)
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	call.Class = p.parseIdentifier()
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	call.Arguments = p.parseCallArguments()
	return call
}
