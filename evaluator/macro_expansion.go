package evaluator

import (
	"github.com/akyrey/monkey-programming-language/ast"
	"github.com/akyrey/monkey-programming-language/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	// Find macro definitions (only top-level, we don't check the child nodes)
	// TODO: walk down the AST and find inner macros
	for i, statement := range program.Statements {
		if isMacroDefined(statement) {
			addMacro(statement, env)
			definitions = append(definitions, i)
		}
	}

	// Remove them from the AST
	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append(
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

// Uses Modify function to recursively walk down the program AST and find calls to macros
// macro callExpressions are evaluated transforming arguments in *object.Quote and extends the environment like
// we do with functions
// It then returns the quoted AST node, replacing the macro call with the result of the evaluation
func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpression)
		evalEnv := extendedMacroEnv(macro, args)

		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

// Just checking if we have a LetStatement with a MacroLiteral
func isMacroDefined(node ast.Statement) bool {
	letStatement, ok := node.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStatement.Value.(*ast.MacroLiteral)
	if !ok {
		return false
	}

	return true
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	// We already check these in isMacroDefined function, so we ignore errors
	letStatement, _ := stmt.(*ast.LetStatement)
	macroLiteral, _ := letStatement.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}

	env.Set(letStatement.Name.Value, macro)
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}

	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}

	for _, a := range exp.Arguments {
		args = append(args, &object.Quote{Node: a})
	}

	return args
}

func extendedMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)

	for i, param := range macro.Parameters {
		extended.Set(param.Value, args[i])
	}

	return extended
}
