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
		return ast.NewBoolean(p.currentToken, &value)
	}

	value = false
	return ast.NewBoolean(p.currentToken, &value)
}

// parse a for expression
func (p *Parser) parseFor() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> por i en rango(10))
		return nil
	}
	condition := p.parseRangeExpression()
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

	body := p.parseBlock()
	return ast.NewFor(token, condition, body)
}

// parse a function expression
func (p *Parser) parseFunction() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	var name *ast.Identifier = nil
	if p.peekToken.Token_type == l.IDENT {
		p.advanceTokens()
		name = p.parseIdentifier().(*ast.Identifier)
	}

	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> funcion {}
		return nil
	}
	parameters := p.parseIdentifiers(l.RPAREN)

	var body *ast.Block
	switch {
	case p.peekToken.Token_type == l.LBRACE:
		p.advanceTokens()
		body = p.parseBlock()

	case p.peekToken.Token_type == l.ARROW:
		p.advanceTokens()
		p.advanceTokens()
		exp := p.parserExpressionStatement()
		body = &ast.Block{Staments: []ast.Stmt{exp}}

	default:
		p.expectedTokenError(l.LBRACE)
		return nil
	}

	return ast.NewFunction(token, name, body, parameters...)
}

// parse a while expression
func (p *Parser) parseWhile() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.LPAREN) {
		// syntax error -> mientras <condition>
		// missing the left paren
		return nil
	}

	p.advanceTokens()
	condition := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		// syntax error -> mientras <condition>) {}
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		// syntax error -> mientras (<condition>) there is no body
		return nil
	}
	body := p.parseBlock()
	return ast.NewWhile(token, condition, body)
}

// parse a identifier expression
func (p *Parser) parseIdentifier() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return ast.NewIdentifier(p.currentToken, p.currentToken.Literal)
}

// parse an if expresion
func (p *Parser) parseIf() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.LPAREN) {
		// syntax error. missing parents
		return nil
	}

	p.advanceTokens()
	condition := p.parseExpression(LOWEST)
	if !p.expepectedToken(l.RPAREN) {
		// syntax error. missing parents
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		// syntax error. missing body
		return nil
	}

	consequence := p.parseBlock()
	var alternative *ast.Block = nil
	p.checkPeekTokenIsNotNil()
	// if we have an else token that means there is an else expression
	if p.peekToken.Token_type == l.ELSE {
		p.advanceTokens()
		if !p.expepectedToken(l.LBRACE) {
			return nil
		}
		alternative = p.parseBlock()
	}

	return ast.NewIf(token, condition, consequence, alternative)
}

// parse a integer expressions
func (p *Parser) parseInteger() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken

	val, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		// the value is not a number. this is very weird to happend
		message := fmt.Sprintf("no se pudo parsear %s como entero", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	return ast.NewInteger(token, &val)
}

// parse a float expression
func (p *Parser) parseFloat() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	val, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		// the value is not a float. this is very weird to happend
		message := fmt.Sprintf("no se pudo parsear %s como flotante", p.currentToken.Literal)
		p.errors = append(p.errors, message)
		return nil
	}

	return ast.NewFloatExp(token, val)
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
	token := p.currentToken
	operator := p.currentToken.Literal

	p.advanceTokens()
	rigth := p.parseExpression(PREFIX)
	return ast.NewPrefix(token, operator, rigth)
}

// parse a string literal
func (p *Parser) parseStringLiteral() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	return ast.NewStringLiteral(p.currentToken, p.currentToken.Literal)
}

// parse a array expression
func (p *Parser) ParseArray() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.LBRACKET) {
		// syntax error -> lista 2,3,4,5
		return nil
	}
	values := p.parseExpressions(l.RBRACKET)
	return ast.NewArray(token, values...)
}

// parse a map expression
func (p *Parser) parseMap() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	keyValues := make([]*ast.KeyValue, 0)
	if !p.expepectedToken(l.LBRACE) {
		// syntax error: missing left brace
		return nil
	}

	p.advanceTokens()
	keyVal := p.parseKeyValues()
	if keyVal != nil {
		keyValues = append(keyValues, keyVal)
	}

	// we loop untile we dont have commas. this means we parse all the key value pairs.
	for p.peekToken.Token_type == l.COMMA {
		p.advanceTokens()
		p.advanceTokens()
		keyVal := p.parseKeyValues()
		if keyVal != nil {
			keyValues = append(keyValues, keyVal)
		}

	}
	p.advanceTokens()
	return ast.NewMapExpression(token, keyValues)
}

// parse a class method
func (p *Parser) parseClassMethod() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	name := p.parseIdentifier().(*ast.Identifier)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	params := p.parseIdentifiers(l.RPAREN)
	var body *ast.Block

	if p.peekToken.Token_type == l.ARROW {
		p.advanceTokens()
		p.advanceTokens()
		exp := p.parseStament()
		body = ast.NewBlock(nil, exp)
	} else {
		if !p.expepectedToken(l.LBRACE) {
			return nil
		}
		body = p.parseBlock()
	}

	p.advanceTokens()
	return ast.NewClassMethodExp(token, name, params, body)
}

// parse a call to instanciate a new class
func (p *Parser) parseClassCall() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	call := ast.NewClassCall(p.currentToken, nil, nil)
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	call.Class = p.parseIdentifier().(*ast.Identifier)
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}

	call.Arguments = p.parseExpressions(l.RPAREN)
	return call
}

// parse an imper statement
func (p *Parser) parseImportStatement() ast.Stmt {
	p.checkCurrentTokenIsNotNil()
	importStmt := ast.NewImportStatement(p.currentToken, nil)
	p.advanceTokens()
	importStmt.Path = p.parseExpression(LOWEST)
	return importStmt
}

func (p *Parser) parseTryExp() ast.Expression {
	try := ast.NewTry(p.currentToken, nil, nil, nil)
	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	try.Try = p.parseBlock()
	if !p.expepectedToken(l.EXCEPT) {
		return nil
	}
	if !p.expepectedToken(l.LPAREN) {
		return nil
	}
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	try.Param = p.parseIdentifier().(*ast.Identifier)
	if !p.expepectedToken(l.RPAREN) {
		return nil
	}

	if !p.expepectedToken(l.LBRACE) {
		return nil
	}

	try.Catch = p.parseBlock()
	return try
}

func (p *Parser) ParseTrhowExp() ast.Expression {
	p.checkCurrentTokenIsNotNil()
	token := p.currentToken
	if !p.expepectedToken(l.IDENT) {
		return nil
	}

	ident := p.parseIdentifier().(*ast.Identifier)
	if ident.Value != "Error" {
		p.errors = append(p.errors, "solo se permite el keyword Error para lanzar una excepcion")
		return nil
	}

	if !p.expepectedToken(l.LPAREN) {
		return nil
	}
	message := p.parseExpressions(l.RPAREN)
	if len(message) != 1 {
		p.errors = append(p.errors, "las excepciones solo pueden recibir un argumento")
		return nil
	}

	return ast.NewThrowExpression(token, message[0])
}
