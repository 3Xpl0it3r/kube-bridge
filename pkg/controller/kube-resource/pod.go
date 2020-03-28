package kube_resource

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/operator/kube-resource"
)

type kubeResourcePodController struct {
	KubeResourceController
}

func NewKubeResourcePodController(clientSet kubernetes.Interface, restConfig *rest.Config, dispatcher controller.IDispatcher)controller.Controller{
	return &kubeResourcePodController{KubeResourceController{
		HookManager: controller.HookManager{},
		clientSet:clientSet,
		operator: kube_resource.NewPodOperator(clientSet, restConfig),
		dispatchor: dispatcher,
	}}
}

func(c *kubeResourcePodController)Run(ctx context.Context)error{
	//informer := informers.NewSharedInformerFactory(c.clientSet, 10 * time.Second).Core().V1().Pods().Informer()
	// define a listwatch for pod which lable contains LABEL_SELECTOR_KUBE_BRIDGE_MODULE

	lw := cache.ListWatch{
		ListFunc: func(options v1.ListOptions) (object runtime.Object, err error) {
			return c.clientSet.CoreV1().Pods(corev1.NamespaceAll).List(v1.ListOptions{
				LabelSelector:       LABEL_SELECTOR_KUBE_BRIDGE_MODULE,
			})
		},
		WatchFunc: func(options v1.ListOptions) (w watch.Interface, err error) {
			return c.clientSet.CoreV1().Pods(corev1.NamespaceAll).Watch(v1.ListOptions{
				LabelSelector:       LABEL_SELECTOR_KUBE_BRIDGE_MODULE,
			})
		},
	}

	podInformer := cache.NewSharedIndexInformer(&lw,&corev1.Pod{}, 0,cache.Indexers{})
	podInformer.AddEventHandlerWithResyncPeriod(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}, 0)

	logging.LogKubeResourceController("pod").Infoln("Pod informer ready to run")
	podInformer.Run(ctx.Done())
	<- ctx.Done()
	return ctx.Err()
}


func(c *kubeResourcePodController)onAdd(object interface{}){

	newObj,ok := object.(*corev1.Pod)
	if !ok {
		logging.LogKubeResourceController("pod").WithField("stage", "onAdd").
			Errorf("Expcet a pod pointer, but get %v\n", object)
		return
	}

	defer func() {
		if err:= c.operator.UpdateOperator(object);err != nil {
			logging.LogKubeResourceController("pod").WithError(err.Error()).Errorf("pod %s state updated failed\n", newObj.Name)
		}
	}()

	if newObj.Status.Phase != "Running" {
		newObj.Labels[KUBE_BRIDGE_MODULE_STATE] = KUBE_BRIDGE_MODULE_UPDATING
		logging.LogKubeResourceController("pod").Debugf("pod :%s state is not  running ,waiting\n", newObj.Name)
		return
	}
	if _,ok := newObj.Labels[KUBE_BRIDGE_MODULE_STATE]; !ok || newObj.Labels[KUBE_BRIDGE_MODULE_STATE]!= KUBE_BRIDGE_MODULE_READY{
		logging.LogKubeResourceController("pod").Debugf("pod %s state is running ,but has not add dns\n", newObj.Name)
		if err := c.operator.AddOperator(object);err != nil {
			logging.LogKubeResourceController("pod").WithError(err.Error()).Errorf("pod %s add dns failed, and retry\n", newObj.Name)
			newObj.Labels[KUBE_BRIDGE_MODULE_STATE] = KUBE_BRIDGE_MODULE_UPDATING
		} else {
			logging.LogKubeResourceController("pod").Infof("pod %s add dns successfully\n", newObj.Name)
			newObj.Labels[KUBE_BRIDGE_MODULE_STATE] = KUBE_BRIDGE_MODULE_READY
		}
		return
	}
	// if not external hook ,then do nothing
	for _, hook := range c.HookManager.GetHooks(){
		// do nothing in here
		hook.OnAdd(object)
	}


}
func(c *kubeResourcePodController)onDelete(object interface{}){
	// if not external hook , do nothing
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnDelete(object)
	}
}
func(c *kubeResourcePodController)onUpdate(oldObj, newObj interface{}){
	// check is pod has configure crosscluster
	logging.LogKubeResourceController("pod").Infof("pod %s is updated,then deliver to onAdd()\n", newObj.(*corev1.Pod).Name)
	c.onAdd(newObj)
	for _, hook := range c.HookManager.GetHooks(){
		hook.OnUpdate(newObj)
	}
}

