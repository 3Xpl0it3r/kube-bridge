package kube_resource

import (
	corev1 "k8s.io/api/core/v1"
	extensionBetav1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"l0calh0st.cn/k8s-bridge/pkg/kberror"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
)

type serviceOperator struct {
	clientSet kubernetes.Interface
	restConfig *rest.Config
	operator Operator
}

func NewServiceOperator(clientSet kubernetes.Interface, restConfig *rest.Config)Operator{
	return &serviceOperator{
		clientSet:clientSet,
		restConfig:restConfig,
		operator:NewIngressOperator(clientSet, restConfig),
	}
}

func(op *serviceOperator)AddOperator(object interface{})kberror.KubeBridgeError{
	logging.LogKubeResourceController("service").WithField("Subsystem", "operator").Debugf("Ready remove service %s\n", object.(*corev1.Service))
	return op.operator.AddOperator(derviedIngressFromService(object.(*corev1.Service)))


}

func(op *serviceOperator)DeleteOperator(object interface{})kberror.KubeBridgeError{
	return op.operator.DeleteOperator(derviedIngressFromService(object.(*corev1.Service)))
}
func(op *serviceOperator)UpdateOperator(object interface{})kberror.KubeBridgeError{
	return nil
}


func derviedIngressFromService(svc *corev1.Service)*extensionBetav1.Ingress{
	ingress := extensionBetav1.Ingress{
		TypeMeta:   v1.TypeMeta{
			Kind: "Ingress",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
				Name:svc.Name,
				Namespace: svc.Namespace,
		},
		Spec:       extensionBetav1.IngressSpec{
					Backend: &extensionBetav1.IngressBackend{
						ServiceName: svc.Name,
						ServicePort: intstr.IntOrString{IntVal:svc.Spec.Ports[0].Port},
					},
		},
		Status:     extensionBetav1.IngressStatus{},
	}
	return &ingress
}