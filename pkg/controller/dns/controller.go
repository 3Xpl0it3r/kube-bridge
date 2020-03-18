package dns

import (
	"context"
	"l0calh0st.cn/k8s-bridge/pkg/controller"
)

type domainNameResolveController struct {

}


func NewDNSController()controller.Controller{
	return &domainNameResolveController{}
}

func(c * domainNameResolveController)Run(ctx context.Context)error{
	return nil
}

func(c *domainNameResolveController)AddHook(hook controller.Hook)error{
	return nil
}
func(c *domainNameResolveController)RemoveHook(hook controller.Hook)error{
	return nil
}
