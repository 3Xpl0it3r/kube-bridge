package kube_resource

import (

	"k8s.io/client-go/kubernetes"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
)


const (
	LABEL_FLAG = "kubebridge"
)

type kubeResourceController struct {
	controller.HookManager
	clientSet kubernetes.Interface
	operator  kube_resource.Operator
}






