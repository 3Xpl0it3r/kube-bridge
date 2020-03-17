package kube_resource

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
)

type kubeResourceServiceController struct {
	kubeResourceController
}

func NewKubeResourceServiceController(clientSet kubernetes.Interface)controller.Controller{
	return &kubeResourceServiceController{kubeResourceController{
		HookManager: controller.HookManager{},
		clientSet:clientSet,
		operator: kube_resource.NewServiceOperator(clientSet),
	}}
}

func(c *kubeResourceServiceController)Run(ctx context.Context)error{
	informer:= informers.NewSharedInformerFactory(c.clientSet, 0).Core().V1().Services().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd	,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})
	informer.Run(ctx.Done())
	<- ctx.Done()
	return ctx.Err()
}


func(c *kubeResourceServiceController)onAdd(object interface{}){
	newObj := object.(*corev1.Service)

	if _, ok := newObj.Labels[LABEL_FLAG] ;!ok  {
		return
	}
	if err := c.operator.AddOperator(object);err != nil {
	}
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnAdd(object)
	}
}

func(c *kubeResourceServiceController)onDelete(obj interface{}){

}
func(c *kubeResourceServiceController)onUpdate(oldObj, newObj interface{}){

}
