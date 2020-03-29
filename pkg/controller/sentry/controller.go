package sentry

import (
	"context"
	"l0calh0st.cn/k8s-bridge/configure"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/logging"
	"l0calh0st.cn/k8s-bridge/pkg/operator/sentry"
)

const (
	WORKQUEUESIZE = 100
)

var globalConig *configure.Config = configure.NewConfig()

type KubeBridgeSyncController struct {
	operator   sentry.Operator
	workQueue  controller.EventsHook
	dispatcher controller.IDispatcher
}

func NewKubeBridgeSentryController(dispatcher controller.IDispatcher)controller.Controller{
	c := &KubeBridgeSyncController{
		dispatcher: dispatcher,
		workQueue: controller.NewEventsHook(WORKQUEUESIZE),
	}
	c.operator = sentry.NewRealSyncOperator(c.workQueue)
	return c
}

func(c *KubeBridgeSyncController)Run(ctx context.Context)error{
	go func() {
		logging.LogSentryController().Infof("Sentry Workqueue is working")
		for true {
			select {
			case record := <- c.workQueue.GetEventChan():
				logging.LogSentryController().Infof("Receive Event Single From Sentry")
				c.dispatcher.Dispatch(record, c)
			case <- ctx.Done():
			}
		}
	}()

	if err := c.operator.Run(ctx);err != nil {
		logging.LogSentryController().WithError(err).Errorf("Run Sentry server Failed")
	}
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
	logging.LogSentryController().Infof("Receive Dispatch From Kubernetes Resource Controller")
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
	err := c.operator.OnAdd(object)
	if err != nil {
		logging.LogSentryController().WithError(err).Errorf("Send Add  to peer cluster Failed")
	} else {
		logging.LogSentryController().Errorf("Send Add  to peer cluster Success")
	}
}

func(c *KubeBridgeSyncController)onUpdate(object interface{}){
	//if err := c.operator.OnUpdate(object);err != nil {
	//	c.dispatcher.Dispatch(controller.Event{
	//		Type:   controller.EventUpdated,
	//		Object: object,
	//	}, c)
	//}
	err := c.operator.OnUpdate(object)
	if err != nil {
		logging.LogSentryController().WithError(err).Errorf("Send Update  to peer cluster Failed")
	} else {
		logging.LogSentryController().Errorf("Send Update  to peer cluster Success")
	}

}
func(c *KubeBridgeSyncController)onDelete(object interface{}){
	//if err := c.operator.OnDelete(object); err != nil {
	//	c.dispatcher.Dispatch(controller.Event{
	//		Type:   controller.EventDeleted,
	//		Object: object,
	//	},c)
	//}
	err := c.operator.OnDelete(object)
	if err != nil {
		logging.LogSentryController().WithError(err).Errorf("Send Delete to peer cluster Failed")
	} else {
		logging.LogSentryController().Errorf("Send Delete to peer cluster Success")
	}
}