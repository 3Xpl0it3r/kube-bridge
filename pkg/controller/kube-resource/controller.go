package kube_resource

import (
	"k8s.io/client-go/kubernetes"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
)


const (
	LABEL_SELECTOR_KUBE_BRIDGE_MODULE = "kubeBridgeType=external"
	KUBE_BRIDGE_MODULE_STATE = "kubeBridgeState"
	KUBE_BRIDGE_MODULE_READY = "READY"
	KUBE_BRIDGE_MODULE_UPDATING = "UPDATING"
)

type kubeResourceController struct {
	controller.HookManager
	clientSet kubernetes.Interface
	operator  kube_resource.Operator
}




