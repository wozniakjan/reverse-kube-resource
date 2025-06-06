/*
Boilerplate 2025 test.
*/

// Code generated by reverse-kube-resource. DO NOT EDIT.

package examples

import (
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// PersistentVolume "pv0003"
	pv0003PersistentVolumeVolumeModeFilesystem corev1.PersistentVolumeMode = "Filesystem"

	pv0003PersistentVolume = corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pv0003",
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: map[corev1.ResourceName]apiresource.Quantity{
				"storage": apiresource.MustParse("5Gi"),
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: "172.17.0.2",
					Path:   "/tmp",
				},
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			PersistentVolumeReclaimPolicy: "Recycle",
			StorageClassName:              "slow",
			MountOptions: []string{
				"hard",
				"nfsvers=4.1",
			},
			VolumeMode: &pv0003PersistentVolumeVolumeModeFilesystem,
		},
	}
)
