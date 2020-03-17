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

type kubeResourcePodController struct {
	kubeResourceController
}

func NewKubeResourcePodController(clientSet kubernetes.Interface)controller.Controller{
	return &kubeResourcePodController{kubeResourceController{
		HookManager: controller.HookManager{},
		clientSet:clientSet,
		operator: kube_resource.NewPodOperator(clientSet),

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
	newObj := object.(*corev1.Pod)
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
	object, _ = object.(*corev1.Pod)
	if err := c.operator.DeleteOperator(object);err != nil {

	}
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnDelete(object)
	}

}
func(c *kubeResourcePodController)onUpdate(oldObj, newObj interface{}){
	oldObj, _ = oldObj.(*corev1.Pod)
	newObj, _ = newObj.(*corev1.Pod)
	if err := c.operator.UpdateOperator(newObj);err != nil {

	}
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnUpdate(newObj)
	}
}

