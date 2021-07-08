package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/kubernetes/scheme"
)

var src = `package xyz
`

type object struct {
	rt runtime.Object
	un *unstructured.Unstructured
}

type imp struct {
	name string
	path string
}

func nameImport(kImportPath string) string {
	s := strings.Split(kImportPath, "/")
	if len(s) > 2 {
		return fmt.Sprintf("%v%v", s[len(s)-2], s[len(s)-1])
	}
	return s[len(s)-1]
}

func missing(un interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}
	if next, ok := un.(map[string]interface{}); ok {
		val, exists := next[path[0]]
		if !exists {
			return true
		}
		return missing(val, path[1:])
	}
	return true
}

func processHelper(name string, o interface{}, un *unstructured.Unstructured, path []string) (imports []imp, lines []string) {
	if missing(un.Object, path) {
		return
	}
	ve := reflect.ValueOf(o)
	te := reflect.TypeOf(o)
	if ve.CanAddr() {
		ve = ve.Elem()
		te = te.Elem()
	}
	ni := nameImport(te.PkgPath())
	if ni != "" {
		imports = append(imports, imp{name: ni, path: te.PkgPath()})
	}

	switch ve.Kind() {
	case reflect.Struct:
		lines = append(lines, fmt.Sprintf("%v: %v.%v{", name, ni, te.Name()))
		for i := 0; i < ve.NumField(); i++ {
			f := ve.Field(i)
			if !f.CanInterface() {
				// skip unexported fields
				continue
			}
			ifc := f.Interface()
			pathHelper := make([]string, len(path))
			copy(pathHelper, path)
			tag := getTag(te.Field(i))
			if tag != "" {
				pathHelper = append(pathHelper, tag)
			}
			name := te.Field(i).Name
			importsHelper, linesHelper := processHelper(name, ifc, un, pathHelper)
			lines = append(lines, linesHelper...)
			imports = append(imports, importsHelper...)
		}
	case reflect.String:
		lines = append(lines, fmt.Sprintf("%v: %q,", name, ve.Interface()))
	default:
		lines = append(lines, fmt.Sprintf("%v: %v,", name, ve.Interface()))
	}

	if ve.Kind() == reflect.Struct {
		lines = append(lines, "},")
	}
	return
}

func process(o interface{}, un *unstructured.Unstructured, path []string) (imports []imp, lines []string) {
	ve := reflect.ValueOf(o).Elem()
	te := reflect.TypeOf(o).Elem()
	ni := nameImport(te.PkgPath())

	varName := fmt.Sprintf("%v%v", un.GetName(), te.Name())
	lines = append(lines, fmt.Sprintf("var %v = %v.%v{", varName, ni, te.Name()))
	imports = append(imports, imp{name: ni, path: te.PkgPath()})
	for i := 0; i < ve.NumField(); i++ {
		f := ve.Field(i)
		if !f.CanInterface() {
			// skip unexported fields
			continue
		}
		ifc := f.Interface()
		if _, ok := ifc.(metav1.TypeMeta); ok {
			// skip type meta as that is schema's job
			continue
		}
		var path []string
		tag := getTag(te.Field(i))
		if tag != "" {
			path = append(path, tag)
		}
		name := te.Field(i).Name
		importsHelper, linesHelper := processHelper(name, ifc, un, path)
		lines = append(lines, linesHelper...)
		imports = append(imports, importsHelper...)
	}
	lines = append(lines, "}")
	return
}

func getTag(f reflect.StructField) string {
	str := string(f.Tag)
	reStr := `json:"([^"]*)"`
	re := regexp.MustCompile(reStr)
	match := re.FindStringSubmatch(str)
	if len(match) == 1 {
		return ""
	}
	s := strings.Split(match[1], ",")
	for _, v := range s[1:] {
		if v == "inline" {
			return ""
		}
	}
	return s[0]
}

func processObjects(obj []object) (imports []imp, allLines []string) {
	_ = v1.Namespace{}
	for _, o := range obj {
		imports2, lines := process(o.rt, o.un, []string{})
		imports = append(imports, imports2...)
		allLines = append(allLines, lines...)
	}
	return
}

func printImports(obj []imp, buf *bytes.Buffer) {
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

func getRuntimeObject(data []byte) runtime.Object {
	// TODO allow adding CRDs to the scheme
	codecs := serializer.NewCodecFactory(scheme.Scheme)
	obj, _, err := codecs.UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("FAILED: %v", err))
	}
	return obj
}

func getUnstructuredObject(data []byte) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	if _, _, err := dec.Decode(data, nil, obj); err != nil {
		panic(err)
	}
	obj.GetName()
	return obj
}

func printLines(lines []string, buf *bytes.Buffer) {
	for _, l := range lines {
		fmt.Fprintln(buf, l)
	}
}

func main() {
	// TODO allow multiple raw manifests as well as helm charts
	data, err := os.ReadFile("./examples/ns.yaml")
	if err != nil {
		panic(err)
	}

	objs := []object{{
		rt: getRuntimeObject(data),
		un: getUnstructuredObject(data),
	}}
	imports, lines := processObjects(objs)
	var buf bytes.Buffer
	printImports(imports, &buf)
	printLines(lines, &buf)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", string(formatted))
}
