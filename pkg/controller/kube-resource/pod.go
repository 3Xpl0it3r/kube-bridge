package kube_resource

import (
	"context"
	"fmt"
	"k8s.io/client-go/rest"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
)

type kubeResourcePodController struct {
	kubeResourceController
}

func NewKubeResourcePodController(clientSet kubernetes.Interface, restConfig *rest.Config)controller.Controller{
	return &kubeResourcePodController{kubeResourceController{
		HookManager: controller.HookManager{},
		clientSet:clientSet,
		operator: kube_resource.NewPodOperator(clientSet, restConfig),
	}}
}

func(c *kubeResourcePodController)Run(ctx context.Context)error{
	informer := informers.NewSharedInformerFactory(c.clientSet, 0).Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd	,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})

	informer.Run(ctx.Done())
	<- ctx.Done()
	return ctx.Err()
}


func(c *kubeResourcePodController)onAdd(object interface{}){
	newObj,ok := object.(*corev1.Pod)
	if !ok {
		return
	}
	fmt.Println(newObj.Name)
	if _, ok := newObj.Labels[LABEL_FLAG] ;!ok  {
		return
	}

	if err := c.operator.AddOperator(object);err != nil {
	}
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnAdd(object)
	}
}
func(c *kubeResourcePodController)onDelete(object interface{}){
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnDelete(object)
	}
}
func(c *kubeResourcePodController)onUpdate(oldObj, newObj interface{}){
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnUpdate(newObj)
	}
}

