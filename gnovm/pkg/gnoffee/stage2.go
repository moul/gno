package gnoffee

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
)

// Stage2 transforms the given AST files based on gnoffee directives
// and returns an AST for a new generated file based on the provided files.
func Stage2(files map[string]*ast.File) (*ast.File, error) {
	return generateFile(files)
}

func generateFile(pkg map[string]*ast.File) (*ast.File, error) {
	exportMapping := make(map[string]string)
	var packageName string

	for _, f := range pkg {
		if packageName == "" && f.Name != nil {
			packageName = f.Name.Name
		}

		// Iterate over all comments in the file.
		for _, commentGroup := range f.Comments {
			for _, comment := range commentGroup.List {
				// Make sure the comment starts a new line.
				if strings.HasPrefix(comment.Text, "//") {
					parts := strings.Fields(comment.Text)
					switch parts[0] {
					case "//gnoffee:export":
						if len(parts) == 4 && parts[2] == "as" {
							k, v := parts[1], parts[3]
							exportMapping[k] = v
						} else {
							return nil, errors.New("invalid gnoffee:export syntax")
						}
					case "//gnoffee:invar":
						return nil, errors.New("unimplemented: invar keyword")
					default:
						return nil, fmt.Errorf("unknown gnoffee keyword: %s", parts[0])
					}
				}
			}
		}
	}

	newFile := &ast.File{
		Name:  &ast.Ident{Name: packageName},
		Decls: make([]ast.Decl, 0),
	}

	// Now, populate the newFile with the necessary declarations based on the exportMapping.
	for k, v := range exportMapping {
		for _, f := range pkg {
			for _, decl := range f.Decls {

				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					iface, ok := typeSpec.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					if typeSpec.Name.Name != v {
						continue
					}

					for _, method := range iface.Methods.List {
						fnDecl := &ast.FuncDecl{
							Name: method.Names[0],
							Doc: &ast.CommentGroup{
								List: []*ast.Comment{
									{
										Text: "\n// This function was generated by gnoffee due to the export directive.",
									},
								},
							},
							Type: method.Type.(*ast.FuncType),
							Body: &ast.BlockStmt{
								List: make([]ast.Stmt, 0),
							},
						}

						callExpr := &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(k),
								Sel: method.Names[0],
							},
							Args: funcTypeToIdentList(method.Type.(*ast.FuncType).Params),
						}

						// Check if the method has return values
						if method.Type.(*ast.FuncType).Results != nil && len(method.Type.(*ast.FuncType).Results.List) > 0 {
							retStmt := &ast.ReturnStmt{
								Results: []ast.Expr{callExpr},
							}
							fnDecl.Body.List = append(fnDecl.Body.List, retStmt)
						} else {
							exprStmt := &ast.ExprStmt{X: callExpr}
							fnDecl.Body.List = append(fnDecl.Body.List, exprStmt)
						}

						newFile.Decls = append(newFile.Decls, fnDecl)
					}
				}
			}
		}
	}

	return newFile, nil
}

func funcTypeToIdentList(fields *ast.FieldList) []ast.Expr {
	var idents []ast.Expr
	for _, field := range fields.List {
		for _, name := range field.Names {
			idents = append(idents, ast.NewIdent(name.Name))
		}
	}
	return idents
}

func findObjectByName(file *ast.File, objectName string) ast.Node {
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name.Name == objectName {
						return s.Type
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.Name == objectName {
							return s.Type
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Name.Name == objectName {
				return d.Type
			}
		}
	}
	return nil
}
