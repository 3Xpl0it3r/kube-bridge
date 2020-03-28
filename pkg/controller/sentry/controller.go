package sentry

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/sentry"
	"sync"
)

const (
	WORKQUEUESIZE = 100
)

type KubeBridgeSyncController struct {
	operator   sentry.Operator
	workQueue  controller.EventsHook
	dispatcher controller.IDispatcher
}

func NewKubeBridgeSyncController(dispatcher controller.IDispatcher)controller.Controller{
	c := &KubeBridgeSyncController{
		dispatcher: dispatcher,
		workQueue: controller.NewEventsHook(WORKQUEUESIZE),
	}
	c.operator = sentry.NewRealSyncOperator(c.workQueue)
	return c
}

func(c *KubeBridgeSyncController)Run(ctx context.Context)error{
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for true {
			select {
			case record := <- c.workQueue.GetEventChan():
				c.dispatcher.Dispatch(record, c)
			case <- ctx.Done():
			}
		}
	}()

	if err := c.operator.Run(ctx);err != nil {
	}
	wg.Wait()
	<- ctx.Done()
	return ctx.Err()
}


func(c *KubeBridgeSyncController)AddHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeSyncController)RemoveHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeSyncController)Dispatch(event controller.Event, controller controller.Controller){
	c.dispatcher.Dispatch(event, c)
}


func(c *KubeBridgeSyncController)Update(event controller.Event){
	switch event.Type {
	case controller.EventAdded:
		c.onAdd(event.Object)
	case controller.EventUpdated:
		c.onUpdate(event.Object)
	case controller.EventDeleted:
		c.onDelete(event.Object)
	}
}

func(c *KubeBridgeSyncController)onAdd(object interface{}){
	//if err := c.operator.OnAdd(object);err != nil {
	//	c.dispatcher.Dispatch(controller.Event{
	//		Type:   controller.EventAdded,
	//		Object: object,
	//	}, c)
	//}
	c.operator.OnAdd(object)
}

func(c *KubeBridgeSyncController)onUpdate(object interface{}){
	//if err := c.operator.OnUpdate(object);err != nil {
	//	c.dispatcher.Dispatch(controller.Event{
	//		Type:   controller.EventUpdated,
	//		Object: object,
	//	}, c)
	//}
	c.operator.OnUpdate(object)

}
func(c *KubeBridgeSyncController)onDelete(object interface{}){
	//if err := c.operator.OnDelete(object); err != nil {
	//	c.dispatcher.Dispatch(controller.Event{
	//		Type:   controller.EventDeleted,
	//		Object: object,
	//	},c)
	//}
	c.operator.OnDelete(object)
}