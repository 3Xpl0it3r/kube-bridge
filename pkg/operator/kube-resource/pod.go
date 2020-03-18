package kube_resource

import (
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

//
func(op *podOperator)AddOperator(object interface{})error{
	addDns := []string{"/bin/sh", "-c", "echo '114.114.114.114' >>  /etc/resolv.conf"}
	return op.executeRemoteCommand(object.(*corev1.Pod), addDns...)
}

func(op *podOperator)DeleteOperator(object interface{})error{
	return nil
}
func(op *podOperator)UpdateOperator(object interface{})error{
	pod := object.(*corev1.Pod)
	_,err := op.clientSet.CoreV1().Pods(pod.Namespace).Update(pod)
	return err
}


func(op *podOperator)executeRemoteCommand(pod *corev1.Pod,cmd ...string)error {
	return PodRemoteCommandExec(op.clientSet, op.restConfig, pod, cmd...)
}