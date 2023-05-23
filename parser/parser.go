package parser

import (
	"github.com/akyrey/monkey-programming-language/ast"
	"github.com/akyrey/monkey-programming-language/lexer"
	"github.com/akyrey/monkey-programming-language/token"
)

type Parser struct {
	l *lexer.Lexer

	// We need to look at the curToken, which is the current token under
	// examination, to decide what to do next, and we also need peekToken for this decision if curToken
	// doesn’t give us enough information.
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    // Construct root node of the AST
	program := &ast.Program{}
    // The statements array is built repeatedly calling nextToken, that advances both curToken and peekToken
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}
