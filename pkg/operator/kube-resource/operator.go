package kube_resource



type Operator interface {
	AddOperator(object interface{})error
	UpdateOperator(object interface{})error
	DeleteOperator(object interface{})error
}

