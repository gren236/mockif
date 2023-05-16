package generator

import (
	"fmt"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"mockif/ifparser"
	"strings"
)

func Generate(parsedPackage ifparser.Package) string {
	res := fmt.Sprintf("package %s\n", parsedPackage.Name)

	res += generateImports(parsedPackage.Imports)

	for _, interf := range parsedPackage.Interfaces {
		res += generateInterfaceImpl(interf)
	}

	// Format the file correctly
	resFormatted, err := format.Source([]byte(res))
	if err != nil {
		log.Fatalln(err)
	}

	return string(resFormatted)
}

func generateImports(imports map[string]string) string {
	res := "import (\n"

	for path, alias := range imports {
		res += fmt.Sprintf("%s \"%s\"\n", alias, path)
	}

	res += ")\n"

	return res
}

func generateInterfaceImpl(interf ifparser.Interf) string {
	res := generateMockStruct(interf)

	for _, method := range interf.Methods {
		res += generateMockMethod(interf.Name, method)
	}

	res += "\n"
	return res
}

func generateMockStruct(interf ifparser.Interf) string {
	caser := cases.Title(language.English, cases.NoLower)

	res := fmt.Sprintf("type mock%s struct {\n", caser.String(interf.Name))

	for _, method := range interf.Methods {
		params := joinParams(method.Args)
		returns := joinParams(method.Returns)

		// Building signature for struct field
		res += fmt.Sprintf(
			"m%s func(%s) (%s)\n",
			caser.String(method.Name),
			strings.Join(params, ", "),
			strings.Join(returns, ", "),
		)
	}

	res += "}\n"

	return res
}

func generateMockMethod(structName string, method ifparser.Method) string {
	caser := cases.Title(language.English, cases.NoLower)

	paramsSign := joinParams(method.Args)
	returns := joinParams(method.Returns)

	receiverName := strings.ToLower(structName[0:1]) + "m"

	// Building signature
	res := fmt.Sprintf(
		"func (%s mock%s) %s(%s) (%s) {\n",
		receiverName,
		caser.String(structName),
		method.Name,
		strings.Join(paramsSign, ", "),
		strings.Join(returns, ", "),
	)

	// Build params slice without types
	var paramsName []string
	for _, param := range method.Args {
		paramsName = append(paramsName, param.Name)
	}

	// Building func body
	// Do not prepend return if function returns nothing
	if len(returns) > 0 {
		res += "return "
	}
	res += fmt.Sprintf("%s.m%s(%s)\n", receiverName, caser.String(method.Name), strings.Join(paramsName, ", "))

	res += "}\n"

	return res
}

func joinParams(params []ifparser.Param) []string {
	var res []string

	for _, arg := range params {
		res = append(res, fmt.Sprintf("%s %s", arg.Name, arg.Type))
	}

	return res
}
