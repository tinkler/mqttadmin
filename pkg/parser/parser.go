/*
parse package and generate interface
*/
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type Field struct {
	Name     string
	Type     string
	Comments []string
}

type Method struct {
	Name     string
	Comments []string
	Args     []Field
	Rets     []Field
}

type Struct struct {
	Name     string
	Comments []string
	Fields   []Field
	Methods  []Method
}

type Package struct {
	Name       string
	Imports    []string
	ImportsMap map[string]string
	Structs    []Struct
}

func toField(t ast.Expr) (f *Field) {
	switch x := t.(type) {
	case *ast.Ident:
		f = &Field{
			Type: x.Name,
		}
	case *ast.StarExpr:
		f = toField(x.X)
		if f != nil {
			f.Type = "*" + f.Type
		}
	case *ast.SelectorExpr:
		f = &Field{
			Type: x.X.(*ast.Ident).Name + "." + x.Sel.Name,
		}
	case *ast.ArrayType:
		f = toField(x.Elt)
		if f != nil {
			f.Type = "[]" + f.Type
		}
	case *ast.MapType:
		k, v := toField(x.Key), toField(x.Value)
		if k == nil || v == nil {
			return
		}
		f.Type = "map[" + toField(x.Key).Type + "]" + toField(x.Value).Type
	}
	return
}

// Path is the path of the package
// modulePath is the name of the module which the package belongs to
func ParsePackage(path string, modulePath string) (*Package, error) {
	var (
		pkg = &Package{
			ImportsMap: map[string]string{},
		}
		pkgImports = map[string]bool{}
	)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	path = modulePath + "/" + strings.TrimPrefix(strings.TrimPrefix(strings.ReplaceAll(path, string(filepath.Separator), "/"), "./"), "../")
	for _, p := range pkgs {
		pkg.Name = p.Name
		pkg.ImportsMap[pkg.Name] = path
		path = strings.TrimSuffix(path, "/"+pkg.Name)
		for fileName, f := range p.Files {
			if strings.HasSuffix(fileName, "_test.go") || strings.HasSuffix(fileName, "_model.go") {
				continue
			}
			ast.Inspect(f, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.ImportSpec:
					impPath := strings.Trim(x.Path.Value, "\"")
					if impLi := strings.LastIndex(impPath, "/"); impLi > 0 {
						if strings.HasSuffix(path, impPath[:impLi]) {
							importName := impPath[impLi+1:]
							if !pkgImports[impPath] {
								pkgImports[impPath] = true
								pkg.Imports = append(pkg.Imports, importName)
								pkg.ImportsMap[importName] = impPath
							}
						}
					}
				case *ast.GenDecl:
					for _, s := range x.Specs {
						switch t := s.(type) {
						case *ast.TypeSpec:
							switch st := t.Type.(type) {
							case *ast.StructType:
								s := Struct{}
								s.Name = t.Name.Name
								if !ast.IsExported(t.Name.Name) {
									continue
								}
								for _, f := range st.Fields.List {
									if len(f.Names) == 0 {
										continue
									}
									field := toField(f.Type)
									if field == nil {
										continue
									}
									if strings.Contains(field.Type, "context.Context") {
										continue
									}
									for _, n := range f.Names {
										field.Name = n.Name
									}
									if field.Name == "" {
										continue
									}
									if f.Comment != nil {
										for _, c := range f.Comment.List {
											field.Comments = append(field.Comments, strings.TrimSpace(strings.TrimPrefix(c.Text, "//")))
										}
									}
									s.Fields = append(s.Fields, *field)
								}

								if x.Doc != nil {
									for _, c := range x.Doc.List {
										s.Comments = append(s.Comments, strings.TrimSpace(strings.TrimPrefix(c.Text, "//")))
									}
								}

								pkg.Structs = append(pkg.Structs, s)
							}
						}
					}
				case *ast.FuncDecl:

					for i, s := range pkg.Structs {
						if x.Recv == nil {
							continue
						}
						if len(x.Recv.List) == 0 {
							continue
						}
						if s.Name == x.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name {
							m := Method{}
							m.Name = x.Name.Name
							if !ast.IsExported(x.Name.Name) {
								continue
							}
							if x.Type.Params == nil || len(x.Type.Params.List) == 0 {
								continue
							}
							var (
								foundContext bool
							)
							for _, f := range x.Type.Params.List {
								rf := toField(f.Type)
								if rf == nil {
									continue
								}
								if strings.Contains(rf.Type, "context.Context") {
									foundContext = true
									continue
								}
								for _, n := range f.Names {
									rf.Name = n.Name
								}
								if rf.Name != "" {
									m.Args = append(m.Args, *rf)
								}
							}
							if !foundContext {
								continue
							}
							var (
								foundError bool
							)
							if x.Type.Results != nil {
								for _, f := range x.Type.Results.List {
									rf := toField(f.Type)
									if rf == nil {
										continue
									}
									if strings.Contains(rf.Type, "error") {
										foundError = true
										continue
									}
									for _, n := range f.Names {
										rf.Name = n.Name
									}
									m.Rets = append(m.Rets, *rf)
								}
							}
							if !foundError {
								continue
							}
							if x.Doc != nil {
								for _, c := range x.Doc.List {
									m.Comments = append(m.Comments, strings.TrimSpace(strings.TrimPrefix(c.Text, "//")))
								}
							}
							pkg.Structs[i].Methods = append(pkg.Structs[i].Methods, m)
						}
					}
				}
				return true
			})
		}
	}
	return pkg, nil
}

func FindStruct(pkg *Package, name string) *Struct {
	for _, s := range pkg.Structs {
		if s.Name == name {
			return &s
		}
	}
	return nil
}
