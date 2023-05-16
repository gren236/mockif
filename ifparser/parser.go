package ifparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"regexp"
	"strings"
)

var noTestFilesRegexp = regexp.MustCompile("^.*_test\\.go$")

type Package struct {
	Name       string
	Imports    map[string]string // path => alias
	Interfaces []Interf
}

type Interf struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name    string
	Args    []Param
	Returns []Param
}

type Param struct {
	Name string
	Type string
}

func ParseDir(src string) Package {
	fset := token.NewFileSet()

	pdir, err := parser.ParseDir(
		fset,
		src,
		func(fi fs.FileInfo) bool {
			// Do not include tests
			if noTestFilesRegexp.MatchString(fi.Name()) {
				return false
			}

			return true
		},
		0,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var parsedInterfaces []Interf
	var pkgName string
	imports := make(map[string]string)

	// Parse AST
	for name, pkg := range pdir {
		pkgName = name

		for _, file := range pkg.Files {
			parsedIfs, parsedImps := parseFile(file)

			appendImports(imports, parsedImps)

			parsedInterfaces = append(parsedInterfaces, parsedIfs...)
		}

		break
	}

	return Package{
		Name:       pkgName,
		Imports:    imports,
		Interfaces: parsedInterfaces,
	}
}

func parseFile(file *ast.File) ([]Interf, map[string]string) {
	var result []Interf
	imports := make(map[string]string)

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				result = append(result, parseInterface(typeSpec.Name.String(), interfaceType))
			}
		}
	}

	// If interfaces found in file, add imports to result
	if len(result) >= 1 {
		for _, imp := range file.Imports {
			name := ""
			if imp.Name != nil {
				name = imp.Name.Name
			}

			imports[strings.Trim(imp.Path.Value, "\"")] = name
		}
	}

	return result, imports
}

func parseInterface(name string, spec *ast.InterfaceType) Interf {
	inter := Interf{Name: name}

	for _, m := range spec.Methods.List {
		mFunc, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		parsedMethod := Method{
			Name: m.Names[0].String(),
		}

		// Parse and add arguments
		parsedMethod.Args = append(parsedMethod.Args, parseParams(mFunc.Params.List)...)

		// Parse and add Returns
		if mFunc.Results != nil {
			parsedMethod.Returns = append(parsedMethod.Returns, parseParams(mFunc.Results.List)...)
		}

		inter.Methods = append(inter.Methods, parsedMethod)
	}

	return inter
}

func parseParams(list []*ast.Field) []Param {
	var resParams []Param

	for _, prm := range list {
		var parsedParam Param

		if prm.Names != nil {
			parsedParam.Name = prm.Names[0].String()
		}

		// TODO Support function as param
		parsedParam.Type = parseType(prm.Type)

		resParams = append(resParams, parsedParam)
	}

	return resParams
}

func parseType(expr ast.Expr) string {
	// TODO Support function as param
	switch v := expr.(type) {
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X.(*ast.Ident).Name, v.Sel.Name)
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", parseType(v.X))
	case *ast.ArrayType:
		var l string
		if lit, ok := v.Len.(*ast.BasicLit); ok {
			l = lit.Value
		}

		return fmt.Sprintf("[%s]%s", l, parseType(v.Elt))
	default:
		return "<unrecognised>"
	}
}

func appendImports(allImps map[string]string, imps map[string]string) {
	for path, alias := range imps {
		allImps[path] = alias
	}
}
