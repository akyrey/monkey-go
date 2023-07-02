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

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {

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
