package controller

import (
	"context"
)

type Controller interface {
	Run(ctx context.Context)error
	AddHook(hook Hook)error
	RemoveHook(hook Hook)error
	//
	Dispatch(object interface{}, controller Controller)
}
