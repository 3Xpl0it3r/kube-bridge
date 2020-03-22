package storage

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
)

type KBStorageController struct {
	sync controller.IDispatcher
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

func(s *KBStorageController)Dispatch(object interface{}, controller controller.Controller){
	s.sync.Dispatch(object, s)
}

func(s *KBStorageController)Update(object interface{}){

}