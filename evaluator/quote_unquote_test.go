package evaluator

import (
	"testing"

	"github.com/akyrey/monkey-programming-language/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)

		if !ok {
			t.Fatalf("expected *object.Quote. Got %T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		// We make sure that the correct ast.Node is wrapped inside of *object.Quote
		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. Got %q. Want %q", quote.Node.String(), tt.expected)
		}
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Integers
		{`quote(unquote(4))`, `4`},
		{`quote(unquote(4 + 4))`, `8`},
		{`quote(8 + unquote(4 + 4))`, `(8 + 8)`},
		{`quote(unquote(4 + 4) + 8)`, `(8 + 8)`},
		// Since we also pass the environment, we can now unquote variables
		{`let foobar = 8; quote(foobar)`, `foobar`},
		{`let foobar = 8; quote(unquote(foobar))`, `8`},
		// Booleans
		{`quote(unquote(true))`, `true`},
		{`quote(unquote(true == false))`, `false`},
		{`quote(unquote(quote(4 + 4)))`, `(4 + 4)`},
		{
			`let quotedInfixExpression = quote(4 + 4);
            quote(unquote(4 + 4) + unquote(quotedInfixExpression))`,
			`(8 + (4 + 4))`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)

		if !ok {
			t.Fatalf("expected *object.Quote. Got %T (%+v)", evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. Got %q. Want %q", quote.Node.String(), tt.expected)
		}
	}
}
