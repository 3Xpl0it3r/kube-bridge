package kube_resource

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"l0calh0st.cn/k8s-bridge/pkg/kberror"
)

type podOperator struct {
	clientSet kubernetes.Interface
	restConfig *rest.Config
}

func NewPodOperator(clientSet kubernetes.Interface, restConfig *rest.Config)Operator{
	return &podOperator{clientSet:clientSet, restConfig:restConfig}
}

//
func(op *podOperator)AddOperator(object interface{})kberror.KubeBridgeError{
	addDns := []string{"/bin/sh", "-c", fmt.Sprintf("echo '%s' >>  /etc/resolv.conf", kubeResourceConfig.Address)}
	return op.executeRemoteCommand(object.(*corev1.Pod), addDns...)
}

func(op *podOperator)DeleteOperator(object interface{})kberror.KubeBridgeError{
	return nil
}
func(op *podOperator)UpdateOperator(object interface{})kberror.KubeBridgeError{
	pod := object.(*corev1.Pod)
	_,err := op.clientSet.CoreV1().Pods(pod.Namespace).Update(pod)
	return kberror.NewKubeOpeatorError().AddError(err)
}


func(op *podOperator)executeRemoteCommand(pod *corev1.Pod,cmd ...string)kberror.KubeBridgeError {
	return PodRemoteCommandExec(op.clientSet, op.restConfig, pod, cmd...)
}