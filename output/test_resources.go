/*
Copyright 2021 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by reverse-kube-resource. DO NOT EDIT.

package csicinder

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilintstr "k8s.io/apimachinery/pkg/util/intstr"
)

var (
	// ServiceAccount "csi-cinder-controller-sa"
	csiCinderControllerSaServiceAccount = corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "csi-cinder-controller-sa",
			Namespace: "kube-system",
		},
	}

	// ServiceAccount "csi-cinder-node-sa"
	csiCinderNodeSaServiceAccount = corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "csi-cinder-node-sa",
			Namespace: "kube-system",
		},
	}

	// StorageClass "csi-cinder-sc-delete"
	csiCinderScDeleteReclaimPolicyDelete      corev1.PersistentVolumeReclaimPolicy = "Delete"
	csiCinderScDeleteAllowVolumeExpansionTrue bool                                 = true

	csiCinderScDeleteStorageClass = storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-cinder-sc-delete",
		},
		Provisioner:          "cinder.csi.openstack.org",
		ReclaimPolicy:        &csiCinderScDeleteReclaimPolicyDelete,
		AllowVolumeExpansion: &csiCinderScDeleteAllowVolumeExpansionTrue,
	}

	// StorageClass "csi-cinder-sc-retain"
	csiCinderScRetainAllowVolumeExpansionTrue bool                                 = true
	csiCinderScRetainReclaimPolicyRetain      corev1.PersistentVolumeReclaimPolicy = "Retain"

	csiCinderScRetainStorageClass = storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-cinder-sc-retain",
		},
		Provisioner:          "cinder.csi.openstack.org",
		ReclaimPolicy:        &csiCinderScRetainReclaimPolicyRetain,
		AllowVolumeExpansion: &csiCinderScRetainAllowVolumeExpansionTrue,
	}

	// ClusterRole "csi-attacher-role"
	csiAttacherRoleClusterRole = rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-attacher-role",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"csinodes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"patch",
				},
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"volumeattachments",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"patch",
				},
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"volumeattachments/status",
				},
			},
		},
	}

	// ClusterRole "csi-provisioner-role"
	csiProvisionerRoleClusterRole = rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-provisioner-role",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"update",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumeclaims",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"storageclasses",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"storage.k8s.io",
				},
				Resources: []string{
					"csinodes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"list",
					"watch",
					"create",
					"update",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"events",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
				},
				APIGroups: []string{
					"snapshot.storage.k8s.io",
				},
				Resources: []string{
					"volumesnapshots",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
				},
				APIGroups: []string{
					"snapshot.storage.k8s.io",
				},
				Resources: []string{
					"volumesnapshotcontents",
				},
			},
		},
	}

	// ClusterRole "csi-snapshotter-role"
	csiSnapshotterRoleClusterRole = rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-snapshotter-role",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"list",
					"watch",
					"create",
					"update",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"events",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"snapshot.storage.k8s.io",
				},
				Resources: []string{
					"volumesnapshotclasses",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"create",
					"get",
					"list",
					"watch",
					"update",
					"delete",
				},
				APIGroups: []string{
					"snapshot.storage.k8s.io",
				},
				Resources: []string{
					"volumesnapshotcontents",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"update",
				},
				APIGroups: []string{
					"snapshot.storage.k8s.io",
				},
				Resources: []string{
					"volumesnapshotcontents/status",
				},
			},
		},
	}

	// ClusterRole "csi-resizer-role"
	csiResizerRoleClusterRole = rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-resizer-role",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumes",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumeclaims",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"persistentvolumeclaims/status",
				},
			},
			rbacv1.PolicyRule{
				Verbs: []string{
					"list",
					"watch",
					"create",
					"update",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"events",
				},
			},
		},
	}

	// ClusterRole "csi-nodeplugin-role"
	csiNodepluginRoleClusterRole = rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-nodeplugin-role",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"update",
					"patch",
				},
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"events",
				},
			},
		},
	}

	// ClusterRoleBinding "csi-attacher-binding"
	csiAttacherBindingClusterRoleBinding = rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-attacher-binding",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-controller-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "csi-attacher-role",
		},
	}

	// ClusterRoleBinding "csi-provisioner-binding"
	csiProvisionerBindingClusterRoleBinding = rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-provisioner-binding",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-controller-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "csi-provisioner-role",
		},
	}

	// ClusterRoleBinding "csi-snapshotter-binding"
	csiSnapshotterBindingClusterRoleBinding = rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-snapshotter-binding",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-controller-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "csi-snapshotter-role",
		},
	}

	// ClusterRoleBinding "csi-resizer-binding"
	csiResizerBindingClusterRoleBinding = rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-resizer-binding",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-controller-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "csi-resizer-role",
		},
	}

	// ClusterRoleBinding "csi-nodeplugin-binding"
	csiNodepluginBindingClusterRoleBinding = rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-nodeplugin-binding",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-node-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "csi-nodeplugin-role",
		},
	}

	// Role "external-resizer-cfg"
	externalResizerCfgRole = rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "external-resizer-cfg",
			Namespace: "kube-system",
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				Verbs: []string{
					"get",
					"watch",
					"list",
					"delete",
					"update",
					"create",
				},
				APIGroups: []string{
					"coordination.k8s.io",
				},
				Resources: []string{
					"leases",
				},
			},
		},
	}

	// RoleBinding "csi-resizer-role-cfg"
	csiResizerRoleCfgRoleBinding = rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "csi-resizer-role-cfg",
			Namespace: "kube-system",
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "csi-cinder-controller-sa",
				Namespace: "kube-system",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "external-resizer-cfg",
		},
	}

	// Service "openstack-cinder-csi-controllerplugin"
	openstackCinderCsiControllerpluginService = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openstack-cinder-csi-controllerplugin",
			Labels: map[string]string{
				"app":       "openstack-cinder-csi",
				"chart":     "openstack-cinder-csi-1.3.8",
				"component": "controllerplugin",
				"heritage":  "Helm",
				"release":   "cinder-csi",
			},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
		},
	}

	// DaemonSet "openstack-cinder-csi-nodeplugin"
	openstackCinderCsiNodepluginMountPropagationBidirectional   corev1.MountPropagationMode = "Bidirectional"
	openstackCinderCsiNodepluginMountPropagationHostToContainer corev1.MountPropagationMode = "HostToContainer"
	openstackCinderCsiNodepluginPrivilegedTrue                  bool                        = true
	openstackCinderCsiNodepluginAllowPrivilegeEscalationTrue    bool                        = true
	openstackCinderCsiNodepluginTypeDirectoryOrCreate           corev1.HostPathType         = "DirectoryOrCreate"
	openstackCinderCsiNodepluginTypeDirectory                   corev1.HostPathType         = "Directory"

	openstackCinderCsiNodepluginDaemonSet = appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openstack-cinder-csi-nodeplugin",
			Labels: map[string]string{
				"app":       "openstack-cinder-csi",
				"chart":     "openstack-cinder-csi-1.3.8",
				"component": "nodeplugin",
				"heritage":  "Helm",
				"release":   "cinder-csi",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       "openstack-cinder-csi",
					"component": "nodeplugin",
					"release":   "cinder-csi",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "openstack-cinder-csi",
						"chart":     "openstack-cinder-csi-1.3.8",
						"component": "nodeplugin",
						"heritage":  "Helm",
						"release":   "cinder-csi",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "socket-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet/plugins/cinder.csi.openstack.org",
									Type: &openstackCinderCsiNodepluginTypeDirectoryOrCreate,
								},
							},
						},
						corev1.Volume{
							Name: "registration-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet/plugins_registry/",
									Type: &openstackCinderCsiNodepluginTypeDirectory,
								},
							},
						},
						corev1.Volume{
							Name: "kubelet-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/kubelet",
									Type: &openstackCinderCsiNodepluginTypeDirectory,
								},
							},
						},
						corev1.Volume{
							Name: "pods-probe-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/dev",
									Type: &openstackCinderCsiNodepluginTypeDirectory,
								},
							},
						},
						corev1.Volume{
							Name: "cloud-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/kubernetes",
								},
							},
						},
						corev1.Volume{
							Name: "cacert",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/cacert",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "node-driver-registrar",
							Image: "k8s.gcr.io/sig-storage/csi-node-driver-registrar:v1.3.0",
							Args: []string{
								"--csi-address=$(ADDRESS)",
								"--kubelet-registration-path=$(DRIVER_REG_SOCK_PATH)",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/csi/csi.sock",
								},
								corev1.EnvVar{
									Name:  "DRIVER_REG_SOCK_PATH",
									Value: "/var/lib/kubelet/plugins/cinder.csi.openstack.org/csi.sock",
								},
								corev1.EnvVar{
									Name: "KUBE_NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/csi",
								},
								corev1.VolumeMount{
									Name:      "registration-dir",
									MountPath: "/registration",
								},
							},
							Lifecycle: &corev1.Lifecycle{
								PreStop: &corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"rm -rf /registration/cinder.csi.openstack.org /registration/cinder.csi.openstack.org-reg.sock",
										},
									},
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "liveness-probe",
							Image: "k8s.gcr.io/sig-storage/livenessprobe:v2.1.0",
							Args: []string{
								"--csi-address=/csi/csi.sock",
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/csi",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "cinder-csi-plugin",
							Image: "docker.io/k8scloudprovider/cinder-csi-plugin:v1.21.0",
							Args: []string{
								"/bin/cinder-csi-plugin",
								"--nodeid=$(NODE_ID)",
								"--endpoint=$(CSI_ENDPOINT)",
								"--cloud-config=$(CLOUD_CONFIG)",
							},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "healthz",
									ContainerPort: 9808,
									Protocol:      "TCP",
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "NODE_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								corev1.EnvVar{
									Name:  "CSI_ENDPOINT",
									Value: "unix://csi/csi.sock",
								},
								corev1.EnvVar{
									Name:  "CLOUD_CONFIG",
									Value: "/etc/kubernetes/cloud-config",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/csi",
								},
								corev1.VolumeMount{
									Name:             "kubelet-dir",
									MountPath:        "/var/lib/kubelet",
									MountPropagation: &openstackCinderCsiNodepluginMountPropagationBidirectional,
								},
								corev1.VolumeMount{
									Name:             "pods-probe-dir",
									MountPath:        "/dev",
									MountPropagation: &openstackCinderCsiNodepluginMountPropagationHostToContainer,
								},
								corev1.VolumeMount{
									Name:      "cacert",
									ReadOnly:  true,
									MountPath: "/etc/cacert",
								},
								corev1.VolumeMount{
									Name:      "cloud-config",
									ReadOnly:  true,
									MountPath: "/etc/kubernetes",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: utilintstr.IntOrString{
											Type:   1,
											IntVal: 0,
											StrVal: "healthz",
										},
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      10,
								PeriodSeconds:       60,
								FailureThreshold:    5,
							},
							ImagePullPolicy: "IfNotPresent",
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"SYS_ADMIN",
									},
								},
								Privileged:               &openstackCinderCsiNodepluginPrivilegedTrue,
								AllowPrivilegeEscalation: &openstackCinderCsiNodepluginAllowPrivilegeEscalationTrue,
							},
						},
					},
					DeprecatedServiceAccount: "csi-cinder-node-sa",
					HostNetwork:              true,
					Tolerations: []corev1.Toleration{
						corev1.Toleration{
							Operator: "Exists",
						},
					},
				},
			},
		},
	}

	// StatefulSet "openstack-cinder-csi-controllerplugin"
	openstackCinderCsiControllerpluginReplicas1 int32 = 1

	openstackCinderCsiControllerpluginStatefulSet = appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openstack-cinder-csi-controllerplugin",
			Labels: map[string]string{
				"app":       "openstack-cinder-csi",
				"chart":     "openstack-cinder-csi-1.3.8",
				"component": "controllerplugin",
				"heritage":  "Helm",
				"release":   "cinder-csi",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &openstackCinderCsiControllerpluginReplicas1,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":       "openstack-cinder-csi",
					"component": "controllerplugin",
					"release":   "cinder-csi",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "openstack-cinder-csi",
						"chart":     "openstack-cinder-csi-1.3.8",
						"component": "controllerplugin",
						"heritage":  "Helm",
						"release":   "cinder-csi",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "socket-dir",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						corev1.Volume{
							Name: "cloud-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/kubernetes",
								},
							},
						},
						corev1.Volume{
							Name: "cacert",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/etc/cacert",
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "csi-attacher",
							Image: "k8s.gcr.io/sig-storage/csi-attacher:v3.1.0",
							Args: []string{
								"--csi-address=$(ADDRESS)",
								"--timeout=3m",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/var/lib/csi/sockets/pluginproxy/csi.sock",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/var/lib/csi/sockets/pluginproxy/",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "csi-provisioner",
							Image: "k8s.gcr.io/sig-storage/csi-provisioner:v2.1.1",
							Args: []string{
								"--csi-address=$(ADDRESS)",
								"--timeout=3m",
								"--default-fstype=ext4",
								"--feature-gates=Topology=true",
								"--extra-create-metadata",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/var/lib/csi/sockets/pluginproxy/csi.sock",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/var/lib/csi/sockets/pluginproxy/",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "csi-snapshotter",
							Image: "k8s.gcr.io/sig-storage/csi-snapshotter:v2.1.3",
							Args: []string{
								"--csi-address=$(ADDRESS)",
								"--timeout=3m",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/var/lib/csi/sockets/pluginproxy/csi.sock",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/var/lib/csi/sockets/pluginproxy/",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "csi-resizer",
							Image: "k8s.gcr.io/sig-storage/csi-resizer:v1.1.0",
							Args: []string{
								"--csi-address=$(ADDRESS)",
								"--timeout=3m",
								"--handle-volume-inuse-error=false",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/var/lib/csi/sockets/pluginproxy/csi.sock",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/var/lib/csi/sockets/pluginproxy/",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "liveness-probe",
							Image: "k8s.gcr.io/sig-storage/livenessprobe:v2.1.0",
							Args: []string{
								"--csi-address=$(ADDRESS)",
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADDRESS",
									Value: "/var/lib/csi/sockets/pluginproxy/csi.sock",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/var/lib/csi/sockets/pluginproxy/",
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
						corev1.Container{
							Name:  "cinder-csi-plugin",
							Image: "docker.io/k8scloudprovider/cinder-csi-plugin:v1.21.0",
							Args: []string{
								"/bin/cinder-csi-plugin",
								"--nodeid=$(NODE_ID)",
								"--endpoint=$(CSI_ENDPOINT)",
								"--cloud-config=$(CLOUD_CONFIG)",
								"--cluster=$(CLUSTER_NAME)",
							},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "healthz",
									ContainerPort: 9808,
									Protocol:      "TCP",
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "NODE_ID",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								corev1.EnvVar{
									Name:  "CSI_ENDPOINT",
									Value: "unix://csi/csi.sock",
								},
								corev1.EnvVar{
									Name:  "CLOUD_CONFIG",
									Value: "/etc/kubernetes/cloud-config",
								},
								corev1.EnvVar{
									Name:  "CLUSTER_NAME",
									Value: "kubernetes",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "socket-dir",
									MountPath: "/csi",
								},
								corev1.VolumeMount{
									Name:      "cacert",
									ReadOnly:  true,
									MountPath: "/etc/cacert",
								},
								corev1.VolumeMount{
									Name:      "cloud-config",
									ReadOnly:  true,
									MountPath: "/etc/kubernetes",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: utilintstr.IntOrString{
											Type:   1,
											IntVal: 0,
											StrVal: "healthz",
										},
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      10,
								PeriodSeconds:       60,
								FailureThreshold:    5,
							},
							ImagePullPolicy: "IfNotPresent",
						},
					},
					DeprecatedServiceAccount: "csi-cinder-controller-sa",
				},
			},
			ServiceName: "openstack-cinder-csi-controllerplugin",
		},
	}

	// CSIDriver "cinder.csi.openstack.org"
	cinderCsiOpenstackOrgAttachRequiredTrue bool = true
	cinderCsiOpenstackOrgPodInfoOnMountTrue bool = true

	cinderCsiOpenstackOrgCSIDriver = storagev1.CSIDriver{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cinder.csi.openstack.org",
		},
		Spec: storagev1.CSIDriverSpec{
			AttachRequired: &cinderCsiOpenstackOrgAttachRequiredTrue,
			PodInfoOnMount: &cinderCsiOpenstackOrgPodInfoOnMountTrue,
			VolumeLifecycleModes: []storagev1.VolumeLifecycleMode{
				"Persistent",
				"Ephemeral",
			},
		},
	}
)
