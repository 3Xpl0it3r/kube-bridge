package storage

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
)

type KBStorageController struct {
	dispatcher controller.IDispatcher
}

func NewStorageController()controller.Controller{
	return &KBStorageController{}
}

func(c *KBStorageController)Run(ctx context.Context)error{
	return nil
}

func(c *KBStorageController)AddHook(hook controller.Hook)error{
	return nil
}

func(c *KBStorageController)RemoveHook(hook controller.Hook)error{
	return nil
}

func(s *KBStorageController)Dispatch(event controller.Event, controller controller.Controller){
	s.dispatcher.Dispatch(event, s)
}

func(s *KBStorageController)Update(event controller.Event){
	switch event.Type {
	case controller.EventAdded:

	}
}