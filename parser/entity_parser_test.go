package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/falconandy/sqlboiler-addons/model"
	"github.com/falconandy/sqlboiler-addons/parser"
)

func TestModelParser_ParseEntityFile(t *testing.T) {
	entity, err := parser.ParseEntityFile(filepath.Join("..", "_testdata", "api_keys.go.txt"))

	assert.Nil(t, err)
	assert.Equal(t, "api_keys.go.txt", entity.FileName)
	assert.Equal(t, "models", entity.Package)
	assert.Equal(t, "APIKey", entity.Name)
	assert.Equal(t, 6, len(entity.Fields))

	assert.Equal(t, []model.EntityField{
		{Name: "ID", Type: "int"},
		{Name: "Key", Type: "string"},
		{Name: "Secret", Type: "string"},
		{Name: "Comment", Type: "null.String"},
		{Name: "Disabled", Type: "bool"},
		{Name: "CreatedAt", Type: "time.Time"},
	}, entity.Fields)

	assert.Equal(t, []string{
		"// Insert a single record using an executor.",
		"Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error",
	}, entity.InsertSignature)

	assert.Equal(t, []string{
		"// Update uses an executor to update the APIKey.",
		"// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.",
		"Update(exec boil.ContextExecutor, columns boil.Columns) (int64, error)",
	}, entity.UpdateSignature)

	assert.Equal(t, "apiKeyColumnsWithoutDefault", entity.ColumnsWithoutDefaultName)
}
