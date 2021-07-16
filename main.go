package main

import (
	rev "github.com/wozniakjan/reverse-kube-resource/pkg"
)

func main() {
	objs := rev.ReadInput()
	imports, goObjects := rev.ProcessObjects(objs)
	rev.Print(imports, goObjects)
}
