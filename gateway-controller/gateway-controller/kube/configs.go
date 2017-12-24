package kube

import (
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
)

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
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
								},
							},
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