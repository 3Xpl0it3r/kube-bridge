package sync

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/sync"
)

type KubeBridgeSyncController struct {
	operator sync.Operator
	dispatcher controller.IDispatcher
}

func NewKubeBridgeSyncController(dispatcher controller.IDispatcher)controller.Controller{
	return &KubeBridgeSyncController{
		dispatcher: dispatcher,
		operator: sync.NewRealSyncOperator(),
	}
}

func(c *KubeBridgeSyncController)Run(ctx context.Context)error{
	return nil
}


func(c *KubeBridgeSyncController)AddHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeSyncController)RemoveHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeSyncController)Dispatch(object interface{}, controller controller.Controller){
	c.dispatcher.Dispatch(object, c)
}


func(c *KubeBridgeSyncController)Update(object interface{}){
	c.operator.Update(object)
}