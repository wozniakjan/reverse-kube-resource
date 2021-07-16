package main

import (
	"flag"
	"fmt"
	"os"

	rev "github.com/wozniakjan/reverse-kube-resource/pkg"
)

func fail(msg string) {
	fmt.Printf("%v\n", msg)
	flag.Usage()
	os.Exit(1)
}

func main() {
	pkg := flag.String("package", "", "Name of the package the output go files should have")
	src := flag.String("src", "", "Source path to either yaml file or directory containing yaml files")
	flag.Parse()
	if *pkg == "" {
		fail("Missing required flag -package")
	}
	if *src == "" {
		fail("Missing required -src flag")
	}

	objs := rev.ReadInput()
	imports, goObjects := rev.ProcessObjects(objs)
	rev.Print(*pkg, imports, goObjects)
}
