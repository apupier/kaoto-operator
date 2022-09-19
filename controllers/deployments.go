package controllers

import (
	"github.com/kaotoIO/kaoto-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetFrontEndDeployment(p KaotoParams, kaoto v1alpha1.Kaoto) *appsv1.Deployment {
	image := kaoto.Spec.Frontend.Image
	return getDeployment(kaoto.Name, p.FrontendName, kaoto.Namespace, p.FrontendName, image, p.FrontendPort, "default")
}

func GetBackendDeployment(p KaotoParams, kaoto v1alpha1.Kaoto) *appsv1.Deployment {
	image := kaoto.Spec.Backend.Image
	return getDeployment(kaoto.Name, p.BackendName, kaoto.Namespace, p.BackendName, image, p.BackendPort, "kaoto-operator-integrator-sa")
}

func getDeployment(kaotoName, name, namespace, imageName, image string, port int32, saName string) *appsv1.Deployment {
	labels := labelsForKaoto(name, kaotoName)
	replicas := int32(1)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: saName,
					Containers: []corev1.Container{{
						Image: image,
						Name:  imageName,
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
							Name:          "port",
						}},
						ImagePullPolicy: "Always",
					}},
				},
			},
		},
	}
	return dep
}
func labelsForKaoto(app, name string) map[string]string {
	return map[string]string{"app": app, "kaoto_cr": name}
}
