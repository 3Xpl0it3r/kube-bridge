package kube_resource

import "l0calh0st.cn/k8s-bridge/pkg/kberror"

type Operator interface {
	AddOperator(object interface{})kberror.KubeBridgeError
	UpdateOperator(object interface{})kberror.KubeBridgeError
	DeleteOperator(object interface{})kberror.KubeBridgeError
}

