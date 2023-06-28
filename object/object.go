package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/akyrey/monkey-programming-language/ast"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	ARRAY_OBJ        = "ARRAY"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
)

// Representation of every Monkey value
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Every time will encounter an integer literal in our source code, will create an ast.IntegerLiteral
// and then , when evaluating the AST, turn it into a object.Integer saving the value in this struct and passing
// around a reference to this struct
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// This is pretty much the same as the integer
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return s.Value
}

// Same logic as the above integer
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Represents the absence of a value
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return "null"
}

// Just a wrapper around an object, so we can keep track of it and can later decide whether to stop evaluation or not
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// This is a simple for of errors.
// In a real world interpreter we'd attach a stack trace, line and column numbers of its origin
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// We also include an Env property, since functions carry their own environment with them
// Parameters and Body are taken directly from the ast definition
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// We leave open what they can do, we just require that they accept zero or more Object arguments
// and return an Object
type BuiltinFunction func(args ...Object) Object

// This is just a wrapper
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// With this new data type and the builtin functions we have created related to this we can add some interesting
// functions:
// * map function
// >> let a = [1, 2, 3, 4];
// >> let double = fn(x) { x * 2 };
// >> map(a, double);
// [2, 4, 6, 8]
//
//	let map = fn(arr, f) {
//	    let iter = fn(arr, accumulated) {
//	        if (len(arr) == 0) {
//	            accumulated
//	        } else {
//	            iter(rest(arr), push(accumulated, f(first(arr))));
//	        }
//	    };
//
//	    iter(arr, []);
//	};
//
// * reduce function
//
//	let reduce = fn(arr, initial, f) {
//	    let iter = fn(arr, result) {
//	        if (len(arr) == 0) {
//	            result
//	        } else {
//	            iter(rest(arr), f(result, first(arr)));
//	        }
//	    };
//
//	    iter(arr, initial);
//	};
//
// * sum function
// >> sum([1, 2, 3, 4, 5]);
// 15
//
//	let sum = fn(arr) {
//	    reduce(arr, 0, fn(initial, el) { initial + el });
//	};
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
