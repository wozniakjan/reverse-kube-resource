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

var src = "package xyz"

type object struct {
	rt runtime.Object
	un *unstructured.Unstructured
}

type goObject struct {
	lines []string
}

type imp struct {
	name string
	path string
}

type processContext struct {
	// inputs
	path   []string
	un     *unstructured.Unstructured
	o      interface{}
	name   string
	parent reflect.Kind

	// inputs/outputs
	goObjects *[]goObject
	imports   *[]imp
}

func (pc processContext) new(path, name string, o interface{}, parent reflect.Kind) processContext {
	pc2 := processContext{}
	pc2.path = make([]string, len(pc.path))
	copy(pc2.path, pc.path)
	if path != "" {
		pc2.path = append(pc2.path, path)
	}
	pc2.un = pc.un
	pc2.o = o
	pc2.name = name
	pc2.parent = parent
	pc2.goObjects = pc.goObjects
	pc2.imports = pc.imports
	return pc2
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

func processMapKey(pc processContext) {
	if missing(pc.un.Object, pc.path) {
		return
	}
	ve := reflect.ValueOf(pc.o)
	te := reflect.TypeOf(pc.o)
	ptrDeref := ""
	for ve.Kind() == reflect.Ptr {
		ve = ve.Elem()
		te = te.Elem()
		ptrDeref = fmt.Sprintf("&%v", ptrDeref)
	}
	ni := nameImport(te.PkgPath())
	if ni != "" {
		(*pc.imports) = append((*pc.imports), imp{name: ni, path: te.PkgPath()})
	}
	/*if ve.Kind() == reflect.Ptr {
	}*/
	last := &(*pc.goObjects)[len((*pc.goObjects))-1]
	last.lines = append(last.lines, fmt.Sprintf("%v", "key"))
	return
}

func processMapValue(pc processContext) {
	if missing(pc.un.Object, pc.path) {
		return
	}
	ve := reflect.ValueOf(pc.o)
	te := reflect.TypeOf(pc.o)
	ptrDeref := ""
	for ve.Kind() == reflect.Ptr {
		ve = ve.Elem()
		te = te.Elem()
		ptrDeref = fmt.Sprintf("&%v", ptrDeref)
	}
	ni := nameImport(te.PkgPath())
	if ni != "" {
		(*pc.imports) = append((*pc.imports), imp{name: ni, path: te.PkgPath()})
	}
	last := &(*pc.goObjects)[len((*pc.goObjects))-1]
	last.lines = append(last.lines, fmt.Sprintf("%v", "val"))
	return
}

func processHelper(pc processContext) {
	if missing(pc.un.Object, pc.path) {
		return
	}
	ve := reflect.ValueOf(pc.o)
	te := reflect.TypeOf(pc.o)
	ptrDeref := ""
	for ve.Kind() == reflect.Ptr {
		ve = ve.Elem()
		te = te.Elem()
		ptrDeref = fmt.Sprintf("&%v", ptrDeref)
	}

	ni := nameImport(te.PkgPath())
	if ni != "" {
		(*pc.imports) = append((*pc.imports), imp{name: ni, path: te.PkgPath()})
	}

	ltype := ""
	if pc.parent != reflect.Slice {
		ltype = fmt.Sprintf("%v: ", pc.name)
	}

	switch kind := ve.Kind(); kind {
	case reflect.Struct:
		last := &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v%v%v.%v{", ltype, ptrDeref, ni, te.Name()))
		for i := 0; i < ve.NumField(); i++ {
			f := ve.Field(i)
			if !f.CanInterface() {
				// skip unexported fields
				continue
			}
			tag := getTag(te.Field(i))
			name := te.Field(i).Name
			pc2 := pc.new(tag, name, f.Interface(), te.Kind())
			processHelper(pc2)
		}
		last = &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, "},")
	case reflect.String:
		last := &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v%q,", ltype, ve.Interface()))
	case reflect.Map:
		valElem := te.Elem()
		valNi := ""
		if ni := nameImport(valElem.PkgPath()); ni != "" {
			valNi = fmt.Sprintf("%v.", ni)
		}
		keyElem := te.Key()
		keyNi := ""
		if ni := nameImport(keyElem.PkgPath()); ni != "" {
			keyNi = fmt.Sprintf("%v.", ni)
		}
		last := &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v: map[%v%v]%v%v{", pc.name, keyNi, keyElem.Name(), valNi, valElem.Name()))
		for _, key := range ve.MapKeys() {
			last := &(*pc.goObjects)[len((*pc.goObjects))-1]
			last.lines = append(last.lines, "{")
			val := ve.MapIndex(key)
			pcKey := pc.new("", pc.name, key.Interface(), te.Kind())
			last = &(*pc.goObjects)[len((*pc.goObjects))-1]
			processMapKey(pcKey)
			last.lines[len(last.lines)-1] = last.lines[len(last.lines)-1] + ":"
			pcVal := pc.new("", pc.name, val.Interface(), te.Kind())
			processMapValue(pcVal)
			last.lines[len(last.lines)-1] = last.lines[len(last.lines)-1] + ","
			last.lines = append(last.lines, "},")
		}
		last = &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, "},")
	case reflect.Slice:
		sliceElem := te.Elem()
		sliceNi := ""
		if ni := nameImport(sliceElem.PkgPath()); ni != "" {
			sliceNi = fmt.Sprintf("%v.", ni)
		}
		last := &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v[]%v%v{", ltype, sliceNi, sliceElem.Name()))
		for i := 0; i < ve.Len(); i++ {
			index := ve.Index(i)
			pc2 := pc.new("", pc.name, index.Interface(), kind)
			processHelper(pc2)
		}
		last = &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, "},")
	default:
		last := &(*pc.goObjects)[len((*pc.goObjects))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v%v,", ltype, ve.Interface()))
	}

	return
}

func process(o interface{}, un *unstructured.Unstructured, path []string) (imports []imp, goObjects []goObject) {
	ve := reflect.ValueOf(o).Elem()
	te := reflect.TypeOf(o).Elem()
	ni := nameImport(te.PkgPath())

	varName := fmt.Sprintf("%v%v", un.GetName(), te.Name())
	imports = append(imports, imp{name: ni, path: te.PkgPath()})
	goObjects = []goObject{goObject{}}
	imports = []imp{}
	last := &goObjects[len(goObjects)-1]
	last.lines = append(last.lines, fmt.Sprintf("var %v = %v.%v{", varName, ni, te.Name()))
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
		pc := processContext{
			name:      name,
			parent:    reflect.Struct,
			o:         ifc,
			un:        un,
			path:      path,
			goObjects: &goObjects,
			imports:   &imports,
		}
		processHelper(pc)
	}
	last = &goObjects[len(goObjects)-1]
	last.lines = append(last.lines, "}")
	return
}

func getTag(f reflect.StructField) string {
	str := string(f.Tag)
	reStr := `json:"([^"]*)"`
	re := regexp.MustCompile(reStr)
	match := re.FindStringSubmatch(str)
	if len(match) <= 1 {
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

func processObjects(obj []object) (allImports []imp, allGoObjects []goObject) {
	_ = v1.Namespace{}
	for _, o := range obj {
		imports, goObjects := process(o.rt, o.un, []string{})
		allImports = append(allImports, imports...)
		allGoObjects = append(allGoObjects, goObjects...)
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

func printLines(goObjects []goObject, buf *bytes.Buffer) {
	for _, o := range goObjects {
		for _, l := range o.lines {
			fmt.Fprintln(buf, l)
		}
	}
}

func main() {
	// TODO allow multiple raw manifests as well as helm charts
	data, err := os.ReadFile("./examples/pv.yaml")
	if err != nil {
		panic(err)
	}

	objs := []object{{
		rt: getRuntimeObject(data),
		un: getUnstructuredObject(data),
	}}
	imports, goObjects := processObjects(objs)
	var buf bytes.Buffer
	printImports(imports, &buf)
	printLines(goObjects, &buf)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("failed to format: %v\n%v\n", err, string(buf.Bytes()))
		panic(err)
	}
	fmt.Printf("%v", string(formatted))
}
