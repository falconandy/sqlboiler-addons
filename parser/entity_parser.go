package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"

	"github.com/falconandy/sqlboiler-addons/model"
)

func ParseEntityFile(filePath string) (*model.Entity, error) {
	entity, err := findEntity(filePath)
	if err != nil {
		return nil, err
	}

	insertSignature, updateSignature, err := findCRUDMethods(filePath, entity.Name)
	if err != nil {
		return nil, err
	}
	entity.InsertSignature = insertSignature
	entity.UpdateSignature = updateSignature

	return entity, nil
}

func findEntity(filePath string) (*model.Entity, error) {
	astFile, err := parser.ParseFile(token.NewFileSet(), filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	fileName := filepath.Base(filePath)
	packageName := astFile.Name.String()

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			entity := &model.Entity{
				FileName: fileName,
				Package:  packageName,
				Name:     typeSpec.Name.String(),
			}

			for _, field := range structType.Fields.List {
				fieldName := field.Names[0].String()
				if fieldName == "R" || fieldName == "L" {
					continue
				}

				fieldType := field.Type
				var typeName string
				switch fieldType := fieldType.(type) {
				case *ast.Ident:
					typeName = fieldType.Name
				case *ast.SelectorExpr:
					typeName = fieldType.X.(*ast.Ident).Name + "." + fieldType.Sel.Name
				default:
					return nil, errors.Errorf("unexpected field type %T", fieldType)
				}
				entity.Fields = append(entity.Fields, model.EntityField{
					Name: fieldName,
					Type: typeName,
				})
			}
			return entity, nil
		}
	}

	return nil, errors.New("entity type not found")
}

func findCRUDMethods(filePath, entityName string) (insertSignature, updateSignature []string, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can't read model file")
	}
	lines := strings.Split(string(data), "\n")
	insertPrefix := fmt.Sprintf("func (o *%s) Insert(", entityName)
	updatePrefix := fmt.Sprintf("func (o *%s) Update(", entityName)
	for i, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		isInsert := strings.HasPrefix(line, insertPrefix)
		isUpdate := strings.HasPrefix(line, updatePrefix)
		if isInsert || isUpdate {
			line = strings.TrimPrefix(line, fmt.Sprintf("func (o *%s) ", entityName))
			line = strings.TrimSuffix(line, " {")

			commentLines := make([]string, 0, 3)
			for j := i - 1; j >= 0; j-- {
				commentLine := strings.TrimRightFunc(lines[j], unicode.IsSpace)
				if !strings.HasPrefix(commentLine, "// ") {
					break
				}
				if strings.HasPrefix(commentLine, "// See boil.Columns.") {
					continue
				}
				commentLines = append(commentLines, commentLine)
			}

			for left, right := 0, len(commentLines)-1; left < right; left, right = left+1, right-1 {
				commentLines[left], commentLines[right] = commentLines[right], commentLines[left]
			}

			if isInsert {
				insertSignature = append(insertSignature, commentLines...)
				insertSignature = append(insertSignature, line)
			} else {
				updateSignature = append(updateSignature, commentLines...)
				updateSignature = append(updateSignature, line)
			}
		}
	}

	if len(insertSignature) == 0 || len(updateSignature) == 0 {
		return nil, nil, errors.New("can't find Insert or Update method")
	}

	return insertSignature, updateSignature, nil
}
