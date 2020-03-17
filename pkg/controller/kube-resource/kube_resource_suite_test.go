package kube_resource_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubeResource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KubeResource Suite")
}
