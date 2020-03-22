package controller

import (
	"l0calh0st.cn/k8s-bridge/pkg/controller/dns"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
	"l0calh0st.cn/k8s-bridge/pkg/controller/storage"
	"l0calh0st.cn/k8s-bridge/pkg/controller/sync"
)

const (
	kube_resource_pod_controller_type = "kube_resource_pod_controller"
	kube_resource_service_controller_type = "kube_resource_service_controller"
	kube_dns_controller_type = "kube_dns_controller"
	kube_storage_controller_type = "kube_storage_controller"
)

type IDispatcher interface {
	Dispatch(object interface{}, controller Controller)
}


type Dispatcher struct {
	controllers map[string]Controller
}

func NewDispatcher()*Dispatcher{
	return &Dispatcher{controllers: map[string]Controller{}}
}

func(s *Dispatcher)RegisterController(controller Controller){
	switch controller.(type) {
	case *dns.KubeBridgeDnsController:
		s.controllers[kube_dns_controller_type] = controller
	case *storage.KBStorageController:
		s.controllers[kube_storage_controller_type] = controller
	default:
	}
}

func(s *Dispatcher)Dispatch(object interface{}, controller Controller){
	if s == nil{return }
	switch controller.(type) {
	case *kube_resource.KubeResourceController:
		s.controllers[kube_storage_controller_type].(*sync.KubeBridgeSyncController).Update(object)
	case *sync.KubeBridgeSyncController:
		s.controllers[kube_dns_controller_type].(*dns.KubeBridgeDnsController).Update(object)
	}
}