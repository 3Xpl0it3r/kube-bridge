package dns

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/dns"
)


type KubeBridgeDnsController struct {
	server dns.Operator
	sync controller.IDispatcher
}


func NewKubeBridgeDnsController(sync controller.IDispatcher)controller.Controller{
	return &KubeBridgeDnsController{
		server: dns.NewRealDnsServer(),
		sync: sync,
	}
}


func(c *KubeBridgeDnsController)Run(ctx context.Context)error{
	c.server.Run(ctx)
	<- ctx.Done()
	return ctx.Err()
}

func(c *KubeBridgeDnsController)AddHook(hook controller.Hook)error{
	return nil
}
func(c *KubeBridgeDnsController)RemoveHook(hook controller.Hook)error{
	return nil
}

func(c *KubeBridgeDnsController)Dispatch(object interface{}, controller controller.Controller){
	c.sync.Dispatch(object, c)
}

func(c *KubeBridgeDnsController)Update(object interface{}){
	c.Update(object)
}