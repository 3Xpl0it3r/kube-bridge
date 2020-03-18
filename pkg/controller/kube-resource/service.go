package kube_resource

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
	"strings"
)

// kubeResourceServiceController is the spec of kubeResourceService
type kubeResourceServiceController struct {
	kubeResourceController
}

// NewKubeResourceServiceControoler(
func NewKubeResourceServiceController(clientSet kubernetes.Interface,restConfig *rest.Config)controller.Controller{
	return &kubeResourceServiceController{kubeResourceController{
		HookManager: controller.HookManager{},
		clientSet:clientSet,
		operator: kube_resource.NewServiceOperator(clientSet, restConfig),
	}}
}
// Run the entrypoint of controller
func(c *kubeResourceServiceController)Run(ctx context.Context)error{
	// define a listwatch
	lw := cache.ListWatch{
		ListFunc: func(options v1.ListOptions) (object runtime.Object, err error) {
			return c.clientSet.CoreV1().Services(corev1.NamespaceAll).List(v1.ListOptions{
				LabelSelector:LABEL_SELECTOR_KUBE_BRIDGE_MODULE,
			})
		},
		WatchFunc: func(options v1.ListOptions) (w watch.Interface, err error) {
			return c.clientSet.CoreV1().Services(corev1.NamespaceAll).Watch(v1.ListOptions{
				LabelSelector:LABEL_SELECTOR_KUBE_BRIDGE_MODULE,
			})
		},
		DisableChunking: false,
	}
	svcInformer := cache.NewSharedIndexInformer(&lw, &corev1.Service{},0,cache.Indexers{})
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd	,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})
	svcInformer.Run(ctx.Done())
	<- ctx.Done()
	return ctx.Err()
}


func(c *kubeResourceServiceController)onAdd(object interface{}){
	newObj, ok := object.(*corev1.Service)
	if !ok {return}

	if _, ok := newObj.Labels[KUBE_BRIDGE_MODULE_STATE];
	ok && strings.Compare(newObj.Labels[KUBE_BRIDGE_MODULE_STATE], KUBE_BRIDGE_MODULE_READY) == 0{return}

	defer func() {
		if err := c.operator.UpdateOperator(newObj);err != nil {
			logging.LogKubeResourceController("service").WithError(err).Errorf("update service %s state failed\n", newObj.Name)
		}
	}()
	if err := c.operator.AddOperator(object);err != nil {
		newObj.Labels[KUBE_BRIDGE_MODULE_STATE] = KUBE_BRIDGE_MODULE_UPDATING
	} else {
		newObj.Labels[KUBE_BRIDGE_MODULE_STATE] = KUBE_BRIDGE_MODULE_READY
	}


}

func(c *kubeResourceServiceController)onDelete(obj interface{}){

	// clearn ingress
	var ok bool
	svcObj,ok := obj.(*corev1.Service)
	if !ok {
		logging.LogKubeResourceController("service").Errorf("onDelete expect a runtine.Object(service), but get %v\n", obj)
		return
	}
	logging.LogKubeResourceController("service").Debugf("Service %s ready to be delete\n", svcObj.Name)
	if err := c.operator.DeleteOperator(svcObj);err != nil {
		logging.LogKubeResourceController("service").WithField("Event", "Delete").WithError(err).Errorf("clearn svc/ingress resource failed")
	} else {
		logging.LogKubeResourceController("service").WithField("Event", "Delete").Infof("Delete Service/Ingress  %s Successfully", svcObj.Name)
	}
}
func(c *kubeResourceServiceController)onUpdate(oldObj, newObj interface{}){
	var ok bool
	oldSvc,ok := oldObj.(*corev1.Service)
	if !ok {return }

	newSvc, ok := newObj.(*corev1.Service)
	if !ok {return }

	if oldSvc.ResourceVersion == newSvc.ResourceVersion {return }

	c.onAdd(newObj)

}
