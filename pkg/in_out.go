package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/tools/go/ast/astutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
)

type crd struct {
	schema.GroupVersion
	runtime.Object
}

func checkFatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func printImports(pkg string, obj []Import, buf *bytes.Buffer) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fmt.Sprintf("package %v", pkg), 0)
	checkFatal(err)
	for _, o := range obj {
		astutil.AddNamedImport(fset, f, o.name, o.path)
	}

	ast.SortImports(fset, f)
	printer.Fprint(buf, fset, f)
}

func printLines(rawVars []RawVar, buf *bytes.Buffer, kubermatic bool) {
	if len(rawVars) == 0 {
		return
	}
	single := false
	if !kubermatic {
		if len(rawVars) == 1 && len(rawVars[0].helpers) == 0 {
			single = true
			fmt.Fprintf(buf, "var ")
		} else {
			fmt.Fprintln(buf, "var (")
		}
	}
	for _, o := range rawVars {
		if !single {
			fmt.Fprintf(buf, "// %v %q\n", o.kind, o.name)
		}
		var helpersSorted []string
		for _, l := range o.helpers {
			helpersSorted = append(helpersSorted, l)
		}
		sort.Strings(helpersSorted)
		for _, l := range helpersSorted {
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
	if !single && !kubermatic {
		fmt.Fprintln(buf, ")")
	}
}

func Print(pkg, boilerplate string, imports []Import, vars []RawVar, kubermatic bool) {
	var buf bytes.Buffer
	if boilerplate != "" {
		bytes, err := os.ReadFile(boilerplate)
		checkFatal(err)
		n := time.Now()
		y, _, _ := n.Date()
		str := strings.Replace(string(bytes), "YEAR", fmt.Sprintf("%d", y), 1)
		fmt.Fprintf(&buf, "%v\n", str)
	}
	fmt.Fprintf(&buf, "// Code generated by reverse-kube-resource. DO NOT EDIT.\n\n")
	printImports(pkg, imports, &buf)
	printLines(vars, &buf, kubermatic)
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("%v", buf.String())
		checkFatal(err)
	}
	fmt.Printf("%v", string(formatted))
}

func readFile(path string) (all [][]byte) {
	data, err := os.ReadFile(path)
	checkFatal(err)
	split := strings.Split(string(data), "---")
	for i := range split {
		if len(split[i]) != 0 {
			all = append(all, []byte(split[i]))
		}
	}
	return
}

func read(path string) (all [][]byte) {
	fi, err := os.Stat(path)
	checkFatal(err)
	if fi.IsDir() {
		err := filepath.Walk(path, func(p string, i fs.FileInfo, err error) error {
			pl := strings.ToLower(p)
			if !(strings.HasSuffix(pl, ".yaml") || strings.HasSuffix(pl, ".yml")) || i.IsDir() {
				return nil
			}
			all = append(all, readFile(p)...)
			return nil
		})
		checkFatal(err)
	} else {
		all = readFile(path)
	}
	return all
}

func getGV(pkg *types.Package) *schema.GroupVersion {
	// TODO: implement
	// parse source code for GV registration
	// there should be just one CRD GVK per package
	return &schema.GroupVersion{
		Group:   "gx1",
		Version: "vx1",
	}
}

func getCRD(gv *schema.GroupVersion, obj types.Object) crd {
	//TODO: implement
	return crd{}
}

func crdSchemeInPackage(files []*ast.File, runtimeObjectInterface *types.Interface) ([]crd, error) {
	config := &types.Config{
		Error: func(e error) {
			panic(e)
		},
		Importer: importer.Default(),
	}

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	fset := token.NewFileSet()
	pkg, err := config.Check("genval", fset, files, &info)
	if err != nil {
		return nil, err
	}

	var crds []crd
	gv := getGV(pkg)
	if gv != nil {
		return crds, nil
	}
	scope := pkg.Scope()
	for _, name := range scope.Names() {
		obj := scope.Lookup(name)
		_, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}
		implements := types.Implements(obj.Type(), runtimeObjectInterface)
		if implements {
			crd := getCRD(gv, obj)
			crds = append(crds, crd)
		}
	}
	return crds, nil
}

func appendScheme(crds []crd) error {
	//TODO: implement
	return nil
}

func updateCRDScheme(crdPackage string, runtimeObjectInterface *types.Interface) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, crdPackage, nil, parser.AllErrors)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		files := []*ast.File{}
		for _, f := range pkg.Files {
			files = append(files, f)
		}
		crds, err := crdSchemeInPackage(files, runtimeObjectInterface)
		if err != nil {
			return err
		}
		err = appendScheme(crds)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRuntimeObjectInterface() (*types.Interface, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", runtimeObjectProgram, 0)
	if err != nil {
		return nil, err
	}

	config := &types.Config{
		Importer: importer.Default(),
	}

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	pkg, e := config.Check("genval", fset, []*ast.File{f}, &info)
	if pkg == nil {
		return nil, e
	}
	i := pkg.Scope().Lookup("Object").Type().Underlying().(*types.Interface)
	//	fmt.Println("implements", types.Implements(s, i))
	return i, nil
}

func updateCRDsScheme(crdPackages string) error {
	if crdPackages == "" {
		return nil
	}
	packages := strings.Split(crdPackages, ",")
	runtimeObjectInterface, err := getRuntimeObjectInterface()
	if err != nil {
		return err
	}
	for _, p := range packages {
		if err := updateCRDScheme(p, runtimeObjectInterface); err != nil {
			return err
		}
	}
	return nil
}

func ReadInput(path, crdPackages, namespace string, unstructured bool) (objs []object) {
	d := read(path)
	err := updateCRDsScheme(crdPackages)
	checkFatal(err)
	codecs := serializer.NewCodecFactory(scheme.Scheme)
	for _, data := range d {
		if unstructured {
			objs = append(objs, object{
				rt: getUnstructuredObject(data),
				un: getUnstructuredObject(data),
			})
		} else {
			objs = append(objs, object{
				rt: getRuntimeObject(data, codecs),
				un: getUnstructuredObject(data),
			})
		}
		if namespace != "" {
			m, ok := objs[len(objs)-1].rt.(metav1.Object)
			if ok {
				m.SetNamespace(namespace)
				objs[len(objs)-1].un.SetNamespace(namespace)
			} else {
				un := objs[len(objs)-1].un
				err := fmt.Errorf("failed to set namespace on %v %v/%v", un.GetKind(), un.GetNamespace(), un.GetName())
				checkFatal(err)
			}
		}
	}
	return
}
