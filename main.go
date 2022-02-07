package main

import (
	"flag"
	"fmt"
	"os"

	rev "github.com/wozniakjan/reverse-kube-resource/pkg"
)

func fail(msg string) {
	fmt.Fprintf(os.Stderr, "%v\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {
	pkg := flag.String("package", "", "Name of the package the output go files should have")
	src := flag.String("src", "", "Source path to either yaml file or directory containing yaml files")
	kubermaticInterfaces := flag.Bool("kubermatic-interfaces", false, "Optional flag to output generation of kubermatic interfaces instead of kubernetes API resources.")
	boilerplate := flag.String("go-header-file", "", "File containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year.")
	crdPackages := flag.String("crd-packages", "", "Comma separated list of packages containing CRD definitions")
	onlyMeta := flag.Bool("only-meta", false, "Generate only metadata, don't go to spec")
	namespace := flag.String("namespace", "", "Overwrite namespace")
	unstructured := flag.Bool("unstructured", false, "Generate objects of type unstructured.Unstructured{} instead of the underlying types")
	flag.Parse()
	if *pkg == "" {
		fail("Missing required flag -package")
	}
	if *src == "" {
		fail("Missing required flag -src")
	}

	objs := rev.ReadInput(*src, *crdPackages, *namespace, *unstructured)
	imports, goObjects := rev.ProcessObjects(objs, *kubermaticInterfaces, *onlyMeta)
	rev.Print(*pkg, *boilerplate, imports, goObjects, *kubermaticInterfaces)
}
