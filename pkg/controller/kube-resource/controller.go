package kube_resource

import (
	"context"
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

type KubeResourceController struct {
	controller.HookManager
	clientSet kubernetes.Interface
	operator  kube_resource.Operator
	dispatchor controller.IDispatcher
}

func(c *KubeResourceController)Dispatch(event controller.Event,controller controller.Controller){
	c.dispatchor.Dispatch(event, c)
}
func(c *KubeResourceController)Run(ctx context.Context)error{return nil}
func(c *KubeResourceController)AddHook(hook controller.Hook)error{return nil}
func(c *KubeResourceController)RemoveHook(hook controller.Hook)error{return nil}
func(c *KubeResourceController)Update(event controller.Event){}



