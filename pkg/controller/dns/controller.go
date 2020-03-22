package dns

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/dns"
)


type KubeBridgeDnsController struct {
	server dns.Operator
	sync controller.ISynchronize
}


func NewKubeBridgeDnsController(sync controller.ISynchronize)controller.Controller{
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

func(c *KubeBridgeDnsController)Sync(object interface{}, controller controller.Controller){
	c.server.UpdateZone(object)

}

func(c *KubeBridgeDnsController)Update(object interface{}){

}