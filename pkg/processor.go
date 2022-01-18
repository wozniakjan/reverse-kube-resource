package pkg

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"

	v1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

type object struct {
	rt runtime.Object
	un *unstructured.Unstructured
}

type RawVar struct {
	helpers map[string]string
	lines   []string
	name    string
	kind    string
}

type Import struct {
	name string
	path string
}

type pathElement struct {
	name  string
	kind  string
	index int
}

type processingContext struct {
	// inputs
	path   []pathElement
	un     *unstructured.Unstructured
	ro     runtime.Object
	o      interface{}
	name   string
	parent reflect.Kind

	// outputs
	rawVars *[]RawVar
	imports *[]Import
}

func sanitize(name string) string {
	return strcase.ToLowerCamel(name)
}

func (pc processingContext) new(path pathElement, name string, o interface{}, parent reflect.Kind) processingContext {
	pc2 := processingContext{}
	pc2.path = make([]pathElement, len(pc.path))
	copy(pc2.path, pc.path)
	if path.name != "" || path.kind != "map" {
		pc2.path = append(pc2.path, path)
	}
	pc2.rawVars = pc.rawVars
	pc2.imports = pc.imports
	pc2.un = pc.un
	pc2.ro = pc.ro
	pc2.o = o
	pc2.name = name
	pc2.parent = parent
	return pc2
}

func nameImport(kImportPath string) string {
	s := strings.Split(kImportPath, "/")
	if len(s) > 2 {
		return fmt.Sprintf("%v%v", s[len(s)-2], s[len(s)-1])
	}
	return s[len(s)-1]
}

func missing(un interface{}, path []pathElement) bool {
	if len(path) == 0 {
		return false
	}
	if path[0].kind == "map" {
		if next, ok := un.(map[string]interface{}); ok {
			val, exists := next[path[0].name]
			if !exists {
				return true
			}
			return missing(val, path[1:])
		}
	} else if path[0].kind == "slice" {
		if next, ok := un.([]interface{}); ok {
			if len(next) > path[0].index {
				val := next[path[0].index]
				return missing(val, path[1:])
			}
		} else {
		}
	}
	return true
}

func processMapKey(pc processingContext) {
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
		(*pc.imports) = append((*pc.imports), Import{name: ni, path: te.PkgPath()})
	}
	last := &(*pc.rawVars)[len((*pc.rawVars))-1]
	if ve.Kind() == reflect.String {
		last.lines = append(last.lines, fmt.Sprintf("%v%q", ptrDeref, ve.Interface()))
	} else {
		last.lines = append(last.lines, fmt.Sprintf("%v%v", ptrDeref, ve.Interface()))
	}
	return
}

func processMapValue(pc processingContext) {
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
		(*pc.imports) = append((*pc.imports), Import{name: ni, path: te.PkgPath()})
	}
	last := &(*pc.rawVars)[len((*pc.rawVars))-1]
	if ve.Kind() == reflect.String {
		last.lines = append(last.lines, fmt.Sprintf("%v%q", ptrDeref, ve.Interface()))
	} else {
		if ve.CanInterface() {
			ifc := ve.Interface()
			if q, ok := ifc.(apiresource.Quantity); ok {
				last.lines = append(last.lines, fmt.Sprintf("apiresource.MustParse(%q)", q.String()))
			} else {
				last.lines = append(last.lines, fmt.Sprintf("%v%v", ptrDeref, ve.Interface()))
			}
		} else {
			last.lines = append(last.lines, fmt.Sprintf("%v%v", ptrDeref, ve.Interface()))
		}
	}
	return
}

func processKubermatic(pc processingContext) {
	te := reflect.TypeOf(pc.o).Elem()
	ni := nameImport(te.PkgPath())
	if ni != "" {
		(*pc.imports) = append((*pc.imports), Import{name: ni, path: te.PkgPath()})
	}
	varName := sanitize(fmt.Sprintf("%v-%v", pc.un.GetName(), te.Name()))
	typeName := fmt.Sprintf("%v.%v", ni, te.Name())

	last := &(*pc.rawVars)[len((*pc.rawVars))-1]
	last.lines = append(last.lines, fmt.Sprintf("func %vCreator() reconciling.Named%vCreatorGetter {", varName, te.Name()))
	last.lines = append(last.lines, fmt.Sprintf("return func() (string, reconciling.%vCreator) {", te.Name()))
	last.lines = append(last.lines, fmt.Sprintf("t := %v.DeepCopy()", varName))
	last.lines = append(last.lines, fmt.Sprintf("return t.Name, func(o *%v) (*%v, error) {", typeName, typeName))
	last.lines = append(last.lines, fmt.Sprintf("return t, nil"))
	last.lines = append(last.lines, fmt.Sprintf("}"))
	last.lines = append(last.lines, fmt.Sprintf("}"))
	last.lines = append(last.lines, fmt.Sprintf("}"))
}

func sortMapKeys(keys []reflect.Value) []reflect.Value {
	if len(keys) < 2 {
		return keys
	}

	switch keys[0].Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sorted := make([]int, 0, len(keys))
		for _, k := range keys {
			sorted = append(sorted, k.Interface().(int))
		}
		sort.Ints(sorted)
		out := make([]reflect.Value, 0, len(keys))
		for _, i := range sorted {
			out = append(out, reflect.ValueOf(i))
		}
		return out
	case reflect.Float32, reflect.Float64:
		sorted := make([]float64, 0, len(keys))
		for _, k := range keys {
			sorted = append(sorted, k.Interface().(float64))
		}
		sort.Float64s(sorted)
		out := make([]reflect.Value, 0, len(keys))
		for _, i := range sorted {
			out = append(out, reflect.ValueOf(i))
		}
		return out
	case reflect.String:
		sorted := make([]string, 0, len(keys))
		for _, k := range keys {
			sorted = append(sorted, k.Interface().(string))
		}
		sort.Strings(sorted)
		out := make([]reflect.Value, 0, len(keys))
		for _, i := range sorted {
			out = append(out, reflect.ValueOf(i))
		}
		return out
	default:
		return keys
	}
}

func processKubernetes(pc processingContext) {
	if missing(pc.un.Object, pc.path) {
		return
	}
	ve := reflect.ValueOf(pc.o)
	te := reflect.TypeOf(pc.o)
	teo := te
	ptrDeref := ""
	for ve.Kind() == reflect.Ptr {
		ve = ve.Elem()
		te = te.Elem()
		ptrDeref = fmt.Sprintf("&%v", ptrDeref)
	}

	ni := nameImport(te.PkgPath())
	if ni != "" {
		(*pc.imports) = append((*pc.imports), Import{name: ni, path: te.PkgPath()})
	}

	ltype := ""
	if pc.parent != reflect.Slice {
		ltype = fmt.Sprintf("%v: ", pc.name)
	}

	teType := te.Name()
	if ni != "" {
		teType = fmt.Sprintf("%v.%v", ni, teType)
	}

	switch kind := ve.Kind(); kind {
	case reflect.Struct:
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v%v%v{", ltype, ptrDeref, teType))
		for i := 0; i < ve.NumField(); i++ {
			f := ve.Field(i)
			if !f.CanInterface() {
				// skip unexported fields
				continue
			}
			tag := getTag(te.Field(i))
			name := te.Field(i).Name
			pc2 := pc.new(pathElement{name: tag, kind: "map"}, name, f.Interface(), te.Kind())
			processKubernetes(pc2)
		}
		last = &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, "},")
	case reflect.String:
		varName := fmt.Sprintf("%q", ve.Interface())
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		if ptrDeref != "" {
			varName = sanitize(fmt.Sprintf("%v-%v-%v", pc.un.GetName(), pc.name, ve.Interface()))
			last.helpers[varName] = fmt.Sprintf("%v %v = %q", varName, teType, ve.Interface())
		}
		last.lines = append(last.lines, fmt.Sprintf("%v%v%v,", ltype, ptrDeref, varName))
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
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v: map[%v%v]%v%v{", pc.name, keyNi, keyElem.Name(), valNi, valElem.Name()))
		keys := sortMapKeys(ve.MapKeys())
		for _, key := range keys {
			last := &(*pc.rawVars)[len((*pc.rawVars))-1]
			val := ve.MapIndex(key)
			pcKey := pc.new(pathElement{kind: "map"}, pc.name, key.Interface(), te.Kind())
			last = &(*pc.rawVars)[len((*pc.rawVars))-1]
			processMapKey(pcKey)
			last.lines[len(last.lines)-1] = last.lines[len(last.lines)-1] + ":"
			pcVal := pc.new(pathElement{kind: "map"}, pc.name, val.Interface(), te.Kind())
			processMapValue(pcVal)
			last.lines[len(last.lines)-1] = last.lines[len(last.lines)-1] + ","
		}
		last = &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, "},")
	case reflect.Slice:
		sliceElem := te.Elem()
		sliceNi := ""
		if ni := nameImport(sliceElem.PkgPath()); ni != "" {
			sliceNi = fmt.Sprintf("%v.", ni)
		}
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v[]%v%v{", ltype, sliceNi, sliceElem.Name()))
		for i := 0; i < ve.Len(); i++ {
			index := ve.Index(i)
			pc2 := pc.new(pathElement{kind: "slice", index: i}, pc.name, index.Interface(), kind)
			processKubernetes(pc2)
		}
		last = &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, "},")
	case reflect.Invalid:
		// this happens when empty structs are used to initialize some value
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		last.lines = append(last.lines, fmt.Sprintf("%v%v%v.%v{},", ltype, ptrDeref, ni, teo.Elem().Name()))
	default:
		varName := fmt.Sprintf("%v", ve.Interface())
		last := &(*pc.rawVars)[len((*pc.rawVars))-1]
		if ptrDeref != "" {
			varName = sanitize(fmt.Sprintf("%v-%v-%v", pc.un.GetName(), pc.name, ve.Interface()))
			last.helpers[varName] = fmt.Sprintf("%v %v = %v", varName, teType, ve.Interface())
		}
		last.lines = append(last.lines, fmt.Sprintf("%v%v%v,", ltype, ptrDeref, varName))
	}

	return
}

func process(o runtime.Object, un *unstructured.Unstructured, kubermatic, onlyMeta bool) (imports []Import, rawVars []RawVar) {
	ve := reflect.ValueOf(o).Elem()
	te := reflect.TypeOf(o).Elem()
	ni := nameImport(te.PkgPath())

	varName := sanitize(fmt.Sprintf("%v-%v", un.GetName(), te.Name()))
	imports = append(imports, Import{name: ni, path: te.PkgPath()})
	rawVars = []RawVar{RawVar{name: un.GetName(), kind: te.Name(), helpers: make(map[string]string)}}
	imports = []Import{}
	last := &rawVars[len(rawVars)-1]
	if kubermatic {
		pc := processingContext{
			o:       o,
			un:      un,
			rawVars: &rawVars,
			imports: &imports,
		}
		processKubermatic(pc)
		return
	}
	last.lines = append(last.lines, fmt.Sprintf("%v = %v.%v{", varName, ni, te.Name()))
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
		if _, ok := ifc.(metav1.ObjectMeta); !ok && onlyMeta {
			// cmdline arg wants to generate just object meta
			continue
		}
		var path []pathElement
		tag := getTag(te.Field(i))
		if tag != "" {
			path = append(path, pathElement{name: tag, kind: "map"})
		}
		name := te.Field(i).Name
		pc := processingContext{
			name:    name,
			parent:  reflect.Struct,
			o:       ifc,
			ro:      o,
			un:      un,
			path:    path,
			rawVars: &rawVars,
			imports: &imports,
		}
		processKubernetes(pc)
	}
	last = &rawVars[len(rawVars)-1]
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

func ProcessObjects(obj []object, kubermatic, onlyMeta bool) (allImports []Import, allRawVars []RawVar) {
	_ = v1.Namespace{}
	for _, o := range obj {
		imports, rawVars := process(o.rt, o.un, kubermatic, onlyMeta)
		allImports = append(allImports, imports...)
		allRawVars = append(allRawVars, rawVars...)
	}
	if kubermatic {
		allImports = append(allImports, Import{path: "k8c.io/kubermatic/v2/pkg/resources/reconciling"})
	}
	return
}

func getRuntimeObject(data []byte, codecs serializer.CodecFactory) runtime.Object {
	obj, _, err := codecs.UniversalDeserializer().Decode(data, nil, nil)
	checkFatal(err)
	return obj
}

func getUnstructuredObject(data []byte) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err := dec.Decode(data, nil, obj)
	checkFatal(err)
	obj.GetName()
	return obj
}
