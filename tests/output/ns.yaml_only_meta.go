// Code generated by reverse-kube-resource. DO NOT EDIT.

package examples

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testNamespace = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name: "test",
	},
}
