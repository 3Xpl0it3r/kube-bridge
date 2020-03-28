package controller

import (
	"reflect"
	"strings"
)

const (
	KUBE_RESOURCE_SERVICE_CONTROLLER = "kubeResourceServiceController"
	DNS_CONTROLLER = "KubeBridgeDnsController"
	SENTRY_CONTROLLER = "KubeBridgeSyncController"
)

type IDispatcher interface {
	Dispatch(event Event, controller Controller)
}


type Dispatcher struct {
	controllers map[string]Controller
}

func NewDispatcher()*Dispatcher{
	return &Dispatcher{controllers: map[string]Controller{}}
}

func(s *Dispatcher)RegisterController(controller Controller){
	controllerType := reflect.TypeOf(controller)
	if controllerType == nil {return }
	switch {
	case strings.Contains(controllerType.String(), KUBE_RESOURCE_SERVICE_CONTROLLER):
		s.controllers[KUBE_RESOURCE_SERVICE_CONTROLLER] = controller
	case strings.Contains(controllerType.String(), DNS_CONTROLLER):
		s.controllers[DNS_CONTROLLER] = controller
	case strings.Contains(controllerType.String(), SENTRY_CONTROLLER):
		s.controllers[SENTRY_CONTROLLER] = controller
	default:
	}
}

func(s *Dispatcher)Dispatch(event Event, controller Controller){
	if s == nil{return }
	controllerType := reflect.TypeOf(controller)
	if controllerType == nil {return }

	switch {
	case strings.Contains(controllerType.String(), KUBE_RESOURCE_SERVICE_CONTROLLER):
		s.controllers[SENTRY_CONTROLLER].Update(event)
	case strings.Contains(controllerType.String(), SENTRY_CONTROLLER):
		s.controllers[DNS_CONTROLLER].Update(event)
	default:
	}


}