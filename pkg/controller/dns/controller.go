package dns

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
	"l0calh0st.cn/k8s-bridge/pkg/operator/dns"
)





type kubeBridgeDnsController struct {
	server dns.Operator
}


func NewKubeBridgeDnsController()controller.Controller{
	return &kubeBridgeDnsController{server:nil}
}


func(c *kubeBridgeDnsController)Run(ctx context.Context)error{
	c.server.Run()
	return nil
}

func(c *kubeBridgeDnsController)AddHook(hook controller.Hook)error{
	return nil
}
func(c *kubeBridgeDnsController)RemoveHook(hook controller.Hook)error{
	return nil
}