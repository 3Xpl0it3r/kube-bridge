package kube_resource_test

import (
	. "github.com/onsi/ginkgo"
	. "l0calh0st.cn/k8s-bridge/pkg/controller/kube-resource"
)

var _ = Describe("Service", func() {
	var (
		svcControllerExample = NewKubeResourceServiceController(nil)
	)
	BeforeEach("Run kube Service Controller", )

	Describe("Ready Controller", func() {
		svcControllerExample.Run(nil)
	})

})
