package kube_resource

import "k8s.io/client-go/kubernetes"

type serviceOperator struct {
	clientSet kubernetes.Interface
}

func NewServiceOperator(clientSet kubernetes.Interface)Operator{
	return &serviceOperator{clientSet:clientSet}
}

func(op *serviceOperator)AddOperator(object interface{})error{
	return nil
}

func(op *serviceOperator)DeleteOperator(object interface{})error{
	return nil
}
func(op *serviceOperator)UpdateOperator(object interface{})error{
	return nil
}

