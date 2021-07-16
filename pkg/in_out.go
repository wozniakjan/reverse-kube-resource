package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func printImports(obj []Import, buf *bytes.Buffer) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	for _, o := range obj {
		astutil.AddNamedImport(fset, f, o.name, o.path)
	}

	ast.SortImports(fset, f)
	printer.Fprint(buf, fset, f)
}

func printLines(rawVars []RawVar, buf *bytes.Buffer) {
	if len(rawVars) == 0 {
		return
	}
	single := false
	if len(rawVars) == 1 && len(rawVars[0].helpers) == 0 {
		single = true
		fmt.Fprintf(buf, "var ")
	} else {
		fmt.Fprintln(buf, "var (")
	}
	for _, o := range rawVars {
		if !single {
			fmt.Fprintf(buf, "// %v %q\n", o.kind, o.name)
		}
		for _, l := range o.helpers {
			fmt.Fprintln(buf, l)
		}
		if len(o.helpers) != 0 {
			fmt.Fprintln(buf, "")
		}
		for _, l := range o.lines {
			fmt.Fprintln(buf, l)
		}
		fmt.Fprintln(buf, "")
	}
	if !single {
		fmt.Fprintln(buf, ")")
	}
}

func Print(imports []Import, vars []RawVar) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "// Code generated. DO NOT EDIT!\n\n")
	printImports(imports, &buf)
	printLines(vars, &buf)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("%v", string(buf.Bytes()))
		panic(err)
	}
	fmt.Printf("%v", string(formatted))
}

func read() [][]byte {
	path := "./examples/openstack-cinder-csi.yaml"
	var all [][]byte
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	split := strings.Split(string(data), "---")
	for i, _ := range split {
		if len(split[i]) != 0 {
			all = append(all, []byte(split[i]))
		}
	}
	return all
}

func ReadInput() (objs []object) {
	d := read()
	for _, data := range d {
		objs = append(objs, object{
			rt: getRuntimeObject(data),
			un: getUnstructuredObject(data),
		})
	}
	return
}
