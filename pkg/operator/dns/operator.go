package dns

import "context"

type Operator interface {
	Run(ctx context.Context)error
	AddZone(object interface{})error
	UpdateZone(object interface{})error
	RemoveZone(object interface{})error
}



