package kube_resource

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type podOperator struct {
	clientSet kubernetes.Interface

}

func NewPodOperator(clientSet kubernetes.Interface)Operator{
	return &podOperator{clientSet:clientSet}
}

func(op *podOperator)AddOperator(object interface{})error{
	fmt.Printf("a new pod is added  pod name %s  namespace :%s", object.(*corev1.Pod).Name, object.(*corev1.Pod).Namespace)
	return nil
}

func(op *podOperator)DeleteOperator(object interface{})error{
	return nil
}
func(op *podOperator)UpdateOperator(object interface{})error{
	return nil
}


func(op *podOperator)executeRemoteCommand(name, namespace string)error {
	return nil
}