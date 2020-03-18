package storage

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
)

type storageController struct {

}

func NewStorageController()controller.Controller{
	return &storageController{}
}

func(c *storageController)Run(ctx context.Context)error{
	return nil
}

func(c *storageController)AddHook(hook controller.Hook)error{
	return nil
}

func(c *storageController)RemoveHook(hook controller.Hook)error{
	return nil
}
