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
	boilerplate := flag.String("go-header-file", "", "File containing boilerplate header text. The string YEAR will be replaced with the current 4-digit year.")
	flag.Parse()
	if *pkg == "" {
		fail("Missing required flag -package")
	}
	if *src == "" {
		fail("Missing required flag -src")
	}

	objs := rev.ReadInput(*src)
	imports, goObjects := rev.ProcessObjects(objs)
	rev.Print(*pkg, *boilerplate, imports, goObjects)
}
