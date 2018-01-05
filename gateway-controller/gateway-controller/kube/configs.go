package kube

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// The default Kubernetes service configuration that parameters are added
// to when an application container is deployed
func DefaultServiceConfig() *apiv1.Service {
	 return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{

		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
				},
			},
			Type: apiv1.ServiceTypeNodePort,
		},
	}
}

// The default Kubernetes deployment configuration that parameters are added
// to when an application container is deployed
func DefaultDeploymentConfig() *appsv1beta2.Deployment {
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{

		},
		Spec: appsv1beta2.DeploymentSpec{
			Selector: &metav1.LabelSelector{

			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{

				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Resources:apiv1.ResourceRequirements{
								Limits: apiv1.ResourceList{
									 apiv1.ResourceCPU : resource.MustParse("250m"),
								},
								Requests:apiv1.ResourceList{
									apiv1.ResourceCPU : resource.MustParse("50m"),
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
								},
							},
							// sleep before stopping to give the Ingress controller
							// enough time to detect the new set of deployments
							// when doing zero downtime roll outs.
							Lifecycle: &apiv1.Lifecycle{
								PreStop: &apiv1.Handler{
									Exec: &apiv1.ExecAction{
										Command: []string{"sleep", "15"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// The default Kubernetes autoscaler configuration that parameters are added
// to when an application container is deployed
func DefaultAutoscalerConfig() *autoscalingv2beta1.HorizontalPodAutoscaler {
	return &autoscalingv2beta1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{

		},
		Spec: autoscalingv2beta1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta1.CrossVersionObjectReference{
				Kind: "Deployment",
			},
		},
	}
}

// The default Kubernetes Ingress configuration that parameters are added
// to when parsing an endpoint configuration file
//
// The file specification can be found here:
// https://github.com/AllanKerr/Express/blob/master/docs/gateway/endpoints-file.md
func DefaultIngressConfig() *extensionsv1beta1.Ingress {
	return &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{

		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue {},
					},
				},
			},
		},
	}
}