package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/falconandy/sqlboiler-addons/model"
	"github.com/falconandy/sqlboiler-addons/parser"
)

func TestParseFunctionSignature_Empty(t *testing.T) {
	rawSignature := "f()"
	signature, ok := parser.ParseFunctionSignature(rawSignature)

	assert.True(t, ok)
	assert.Equal(t, "f", signature.Name)
	assert.Empty(t, signature.Parameters)
	assert.Equal(t, "", signature.Return)
}

func TestParseFunctionSignature_Simple(t *testing.T) {
	rawSignature := "doIt(a int,  b, c string, d  *[]int) (int,  string)"
	signature, ok := parser.ParseFunctionSignature(rawSignature)

	assert.True(t, ok)
	assert.Equal(t, "doIt", signature.Name)
	assert.Equal(t, []model.FunctionParameter{
		{
			Names: []string{"a"},
			Type:  "int",
		},
		{
			Names: []string{"b", "c"},
			Type:  "string",
		},
		{
			Names: []string{"d"},
			Type:  "*[]int",
		},
	}, signature.Parameters)
	assert.Equal(t, "(int,  string)", signature.Return)
}
