package main

import (
	"fmt"
	"go/ast"
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
	fmt.Println("nameImport", kImportPath)
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
		return missing(next, path[1:])
	}
	return true
}

func printFields(o interface{}, un *unstructured.Unstructured, path []string) (imports []imp, lines []string) {
	if missing(un, path) {
		return
	}
	v := reflect.ValueOf(o).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}
		ifc := f.Interface()
		if _, ok := ifc.(metav1.TypeMeta); ok {
			continue
		}
		fmt.Printf("%v [%T]\n", i, ifc)
		t := reflect.TypeOf(ifc)
		k := t.Kind()
		e := reflect.TypeOf(o).Elem()
		ef := e.Field(i)
		if k == reflect.Struct {
			var path2 []string
			copy(path, path2)
			tag := getTag(ef)
			if tag != "" {
				path2 = append(path2, tag)
			}
			imports2, _ := printFields(t, un, path2)
			imports = append(imports, imports2...)
		}
		fe := reflect.TypeOf(ifc)
		name := nameImport(fe.PkgPath())
		imports = append(imports, imp{name: name, path: fe.PkgPath()})
	}
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

func printObjects(obj []object) (imports []imp) {
	_ = v1.Namespace{}
	for _, o := range obj {
		imports2, _ := printFields(o.rt, o.un, []string{})
		imports = append(imports, imports2...)
	}
	return
}

func printImports(obj []imp) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	for _, o := range obj {
		astutil.AddNamedImport(fset, f, o.name, o.path)
	}

	ast.SortImports(fset, f)
	printer.Fprint(os.Stdout, fset, f)
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
	imports := printObjects(objs)
	printImports(imports)
	fmt.Println("done")
}
