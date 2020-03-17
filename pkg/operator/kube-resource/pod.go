package kube_resource

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type podOperator struct {
	clientSet kubernetes.Interface
	restConfig *rest.Config
}

func NewPodOperator(clientSet kubernetes.Interface, restConfig *rest.Config)Operator{
	return &podOperator{clientSet:clientSet, restConfig:restConfig}
}

func(op *podOperator)AddOperator(object interface{})error{

	pod := object.(*corev1.Pod)

	if pod.Status.Phase == "Running" {
		err := PodRemoteCommandExec(op.clientSet, op.restConfig, object.(*corev1.Pod))
		if err != nil {
			fmt.Println("remote exec faled  ", err)
		}
		pod.Labels["configed"] = "yes"
		op.clientSet.CoreV1().Pods(pod.Namespace).Update(pod)
	} else {
		fmt.Println("pod is in , wait", object.(*corev1.Pod).Status.Phase)
	}


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