package controller

import (
	"l0calh0st.cn/k8s-bridge/pkg/controller/dns"
	kube_resource "l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
	"l0calh0st.cn/k8s-bridge/pkg/controller/storage"
)

const (
	kube_resource_pod_controller_type = "kube_resource_pod_controller"
	kube_resource_service_controller_type = "kube_resource_service_controller"
	kube_dns_controller_type = "kube_dns_controller"
	kube_storage_controller_type = "kube_storage_controller"
)

type ISynchronize interface {
	Sync(object interface{}, controller Controller)
}


type Synchronize struct {
	controllers map[string]Controller
}

func NewSynchronize()*Synchronize{
	return &Synchronize{controllers: map[string]Controller{}}
}

func(s *Synchronize)RegisterController(controller Controller){
	switch controller.(type) {
	case *dns.KubeBridgeDnsController:
		s.controllers[kube_dns_controller_type] = controller
	case *storage.KBStorageController:
		s.controllers[kube_storage_controller_type] = controller
	default:
	}
}

func(s *Synchronize)Sync(object interface{}, controller Controller){
	if s == nil{return }
	switch controller.(type) {
	case *kube_resource.KubeResourceController:
		s.controllers[kube_storage_controller_type].(*storage.KBStorageController).Update(object)
	case *dns.KubeBridgeDnsController:
		s.controllers[kube_dns_controller_type].(*dns.KubeBridgeDnsController).Update(object)
	}
}