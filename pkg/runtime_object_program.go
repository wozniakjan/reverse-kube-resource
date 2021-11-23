package pkg

var runtimeObjectProgram = `package runtime_object_program

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() runtime.Object
}`
