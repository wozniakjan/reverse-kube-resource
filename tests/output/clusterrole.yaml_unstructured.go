// Code generated by reverse-kube-resource. DO NOT EDIT.

package examples

import v1unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

var clusterRoleUnstructuredClusterRole = v1unstructured.Unstructured{
	Object: map[string]interface{}{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "ClusterRole",
		"metadata": map[string]interface{}{
			"name": "cluster-role",
		},
		"rules": []interface{}{
			map[string]interface{}{
				"apiGroups": []interface{}{
					"",
				},
				"resources": []interface{}{
					"persistentvolumes",
				},
				"verbs": []interface{}{
					"get", "list", "watch", "patch",
				},
			}, map[string]interface{}{
				"apiGroups": []interface{}{
					"storage.k8s.io",
				},
				"resources": []interface{}{
					"csinodes",
				},
				"verbs": []interface{}{
					"get", "list", "watch",
				},
			}, map[string]interface{}{
				"apiGroups": []interface{}{
					"storage.k8s.io",
				},
				"resources": []interface{}{
					"volumeattachments",
				},
				"verbs": []interface{}{
					"get", "list", "watch", "patch",
				},
			}, map[string]interface{}{
				"apiGroups": []interface{}{
					"storage.k8s.io",
				},
				"resources": []interface{}{
					"volumeattachments/status",
				},
				"verbs": []interface{}{
					"patch",
				},
			},
		},
	},
}
