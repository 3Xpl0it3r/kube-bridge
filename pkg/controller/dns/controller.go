package dns

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/dns"
)


type KubeBridgeDnsController struct {
	server dns.Operator
	dispatcher controller.IDispatcher
}


func NewKubeBridgeDnsController(sync controller.IDispatcher)controller.Controller{
	return &KubeBridgeDnsController{
		server: dns.NewRealDnsServer(),
		dispatcher: sync,
	}
}


func(c *KubeBridgeDnsController)Run(ctx context.Context)error{
	if err := c.server.Run(ctx);err != nil {

	}
	<- ctx.Done()
	return ctx.Err()
}

func(c *KubeBridgeDnsController)AddHook(hook controller.Hook)error{
	return nil
}
func(c *KubeBridgeDnsController)RemoveHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeDnsController)Dispatch(event controller.Event, controller controller.Controller){
	c.dispatcher.Dispatch(event, c)
}

func(c *KubeBridgeDnsController)Update(event controller.Event){
	switch event.Type {
	case controller.EventAdded:
		c.onAdd(event.Object)
	case controller.EventDeleted:
		c.onDelete(event.Object)
	case controller.EventUpdated:
		c.onUpdate(event.Object)
	}
}

func(c *KubeBridgeDnsController)onAdd(object interface{}){
	if err := c.server.AddZone(object);err != nil {

	}
}
func(c *KubeBridgeDnsController)onUpdate(object interface{}){
	if err := c.server.UpdateZone(object);err != nil {

	}
}
func(c *KubeBridgeDnsController)onDelete(object interface{}){
	if err := c.server.RemoveZone(object);err != nil {

	} 
}