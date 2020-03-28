package controller

import (
	"context"
)

type Controller interface {
	Run(ctx context.Context)error
	AddHook(hook Hook)error
	RemoveHook(hook Hook)error
	//
	Dispatch(event Event, controller Controller)
	Update(event Event)
}
