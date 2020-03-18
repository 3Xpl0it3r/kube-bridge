package kube_resource

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	extensionBetav1 "k8s.io/api/extensions/v1beta1"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
)

type ingressOperator struct {
	clientSet kubernetes.Interface
	restConfig *rest.Config
}

func NewIngressOperator(clientSet kubernetes.Interface, restConfig *rest.Config)Operator{
	return &ingressOperator{
		clientSet:  clientSet,
		restConfig: restConfig,
	}
}

func(op *ingressOperator)AddOperator(object interface{})error{
	ig := object.(*extensionBetav1.Ingress)
	_,err := op.clientSet.ExtensionsV1beta1().Ingresses(ig.Namespace).Create(ig)
	return err

}
func(op *ingressOperator)DeleteOperator(object interface{})error{

	ig := object.(*extensionBetav1.Ingress)
	logging.LogKubeResourceController("Ingress").WithField("Event","Delete").Debugf("Ready to delete ingress %s\n", ig.Name)
	err := op.clientSet.ExtensionsV1beta1().Ingresses(ig.Namespace).Delete(ig.Name, &v1.DeleteOptions{})
	return err
}
func(op *ingressOperator)UpdateOperator(object interface{})error{
	return nil
}
